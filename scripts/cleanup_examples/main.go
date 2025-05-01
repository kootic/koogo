// cleanup removes all example code prefixed with "Koo"/"koo" from the boilerplate.
//
// Usage:
//
//	go run scripts/cleanup/main.go [--dry-run]
//
// What it does:
//  1. Deletes all files matching koo_*.go and *_koo_*.sql
//  2. Parses Go files and removes Koo-prefixed declarations (types, funcs, vars, consts)
//  3. Removes struct fields with Koo-prefixed types
//  4. Removes route entries containing "/koo/"
//  5. Cleans up unused imports
//
// The script respects the convention:
//   - Files: koo_*.go, *_koo_*.sql
//   - Types/Interfaces: KooUser, KooPet, KooUserService, etc.
//   - Functions: NewKooUserHandler, KooUserFromDomain, etc.
//   - Struct fields with Koo-prefixed types
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var dryRun bool

func main() {
	flag.BoolVar(&dryRun, "dry-run", false, "Print what would be done without making changes")
	flag.Parse()

	root, err := findProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding project root: %v\n", err)
		os.Exit(1)
	}

	if dryRun {
		fmt.Println("=== DRY RUN MODE - No changes will be made ===")
		fmt.Println()
	}

	fmt.Println("=== Cleaning up Koo-prefixed example code ===")
	fmt.Println()

	// Step 1: Delete koo_*.go files and *_koo_*.sql migrations
	if err := deleteKooFiles(root); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting files: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Process remaining Go files to remove Koo references
	if err := processGoFiles(root); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing Go files: %v\n", err)
		os.Exit(1)
	}

	// Step 3: Run goimports to fix imports
	if !dryRun {
		fmt.Println()
		fmt.Println("Running goimports to fix imports...")

		if err := runGoimports(root); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: goimports failed: %v\n", err)
			fmt.Println("  Please run 'goimports -w ./internal/' manually")
		} else {
			fmt.Println("  ✓ Imports fixed")
		}
	}

	fmt.Println()
	fmt.Println("=== Cleanup complete! ===")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Run 'task atlas:hash' to update migration checksums")
	fmt.Println("  2. Run 'go mod tidy' to clean up dependencies")
	fmt.Println("  3. Run 'go build ./...' to verify the build")
	fmt.Println("  4. If you have existing koo_* tables, drop them or reset the database")
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find go.mod in any parent directory")
		}

		dir = parent
	}
}

func deleteKooFiles(root string) error {
	fmt.Println("Deleting koo_* files...")

	var filesToDelete []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and vendor
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}

			return nil
		}

		base := filepath.Base(path)

		// Match koo_*.go files
		if strings.HasPrefix(base, "koo_") && strings.HasSuffix(base, ".go") {
			filesToDelete = append(filesToDelete, path)
			return nil
		}

		// Match *_koo_*.sql migration files
		if strings.Contains(base, "_koo_") && strings.HasSuffix(base, ".sql") {
			filesToDelete = append(filesToDelete, path)
			return nil
		}

		return nil
	})
	if err != nil {
		return err
	}

	for _, path := range filesToDelete {
		rel, _ := filepath.Rel(root, path)
		if dryRun {
			fmt.Printf("  Would delete: %s\n", rel)
		} else {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to delete %s: %w", path, err)
			}

			fmt.Printf("  ✓ Deleted: %s\n", rel)
		}
	}

	if len(filesToDelete) == 0 {
		fmt.Println("  No koo_* files found")
	}

	return nil
}

func processGoFiles(root string) error {
	fmt.Println()
	fmt.Println("Processing Go files to remove Koo references...")

	// Directories to process (internal only, not pkg which has reusable code)
	dirsToProcess := []string{
		filepath.Join(root, "internal"),
	}

	for _, dir := range dirsToProcess {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if !strings.HasSuffix(path, ".go") {
				return nil
			}

			// Skip test files that aren't koo_ prefixed (those are already deleted)
			base := filepath.Base(path)
			if strings.HasPrefix(base, "koo_") {
				return nil // Already handled by deleteKooFiles
			}

			return processGoFile(root, path)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func processGoFile(root, path string) error {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	modified := false

	// Remove Koo-prefixed declarations
	var newDecls []ast.Decl
	for _, decl := range node.Decls {
		if shouldRemoveDecl(decl) {
			modified = true
			continue
		}

		// Check for struct fields to remove
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						if removeKooFields(structType) {
							modified = true
						}
					}
				}
			}
		}

		// Check for function bodies that need modification
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Body != nil {
				if removeKooStatements(funcDecl.Body) {
					modified = true
				}

				if removeKooRoutes(funcDecl.Body) {
					modified = true
				}

				if removeKooStructLitFields(funcDecl.Body) {
					modified = true
				}
			}
		}

		newDecls = append(newDecls, decl)
	}

	node.Decls = newDecls

	// Note: We don't clean up imports here - let goimports handle that

	if !modified {
		return nil
	}

	rel, _ := filepath.Rel(root, path)

	if dryRun {
		fmt.Printf("  Would modify: %s\n", rel)
		return nil
	}

	// Write the modified file
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return fmt.Errorf("failed to format %s: %w", path, err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}

	fmt.Printf("  ✓ Modified: %s\n", rel)

	return nil
}

func shouldRemoveDecl(decl ast.Decl) bool {
	switch d := decl.(type) {
	case *ast.GenDecl:
		// Check type, var, const declarations
		if d.Tok == token.TYPE || d.Tok == token.VAR || d.Tok == token.CONST {
			// Remove if all specs are Koo-prefixed
			allKoo := true
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if !hasKooPrefix(s.Name.Name) {
						allKoo = false
					}
				case *ast.ValueSpec:
					for _, name := range s.Names {
						if !hasKooPrefix(name.Name) {
							allKoo = false
						}
					}
				}
			}

			return allKoo && len(d.Specs) > 0
		}

	case *ast.FuncDecl:
		// Remove Koo-prefixed functions
		if hasKooPrefix(d.Name.Name) {
			return true
		}
		// Remove methods on Koo-prefixed receivers
		if d.Recv != nil && len(d.Recv.List) > 0 {
			recvType := d.Recv.List[0].Type
			if star, ok := recvType.(*ast.StarExpr); ok {
				if ident, ok := star.X.(*ast.Ident); ok {
					if hasKooPrefix(ident.Name) {
						return true
					}
				}
			} else if ident, ok := recvType.(*ast.Ident); ok {
				if hasKooPrefix(ident.Name) {
					return true
				}
			}
		}
	}

	return false
}

func hasKooPrefix(name string) bool {
	return strings.HasPrefix(name, "Koo") || strings.HasPrefix(name, "koo") || strings.Contains(name, "Koo")
}

func removeKooFields(structType *ast.StructType) bool {
	if structType.Fields == nil {
		return false
	}

	modified := false

	var newFields []*ast.Field
	for _, field := range structType.Fields.List {
		if fieldHasKooType(field) {
			modified = true
			continue
		}

		newFields = append(newFields, field)
	}

	structType.Fields.List = newFields

	return modified
}

func fieldHasKooType(field *ast.Field) bool {
	// Check field names
	for _, name := range field.Names {
		if hasKooPrefix(name.Name) {
			return true
		}
	}

	// Check field type
	return typeHasKooPrefix(field.Type)
}

func typeHasKooPrefix(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.Ident:
		return hasKooPrefix(t.Name)
	case *ast.SelectorExpr:
		return hasKooPrefix(t.Sel.Name)
	case *ast.StarExpr:
		return typeHasKooPrefix(t.X)
	case *ast.ArrayType:
		return typeHasKooPrefix(t.Elt)
	}

	return false
}

func removeKooStatements(block *ast.BlockStmt) bool {
	if block == nil {
		return false
	}

	modified := false

	var newStmts []ast.Stmt
	for _, stmt := range block.List {
		if stmtHasKooRef(stmt) {
			modified = true
			continue
		}

		newStmts = append(newStmts, stmt)
	}

	block.List = newStmts

	return modified
}

func stmtHasKooRef(stmt ast.Stmt) bool {
	s, ok := stmt.(*ast.AssignStmt)
	if !ok {
		return false
	}

	// Check for assignments like: userHandler := NewKooUserHandler(...)
	for _, rhs := range s.Rhs {
		if callExpr, ok := rhs.(*ast.CallExpr); ok {
			if ident, ok := callExpr.Fun.(*ast.Ident); ok {
				if hasKooPrefix(ident.Name) {
					return true
				}
			}

			if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if hasKooPrefix(sel.Sel.Name) {
					return true
				}
			}
		}
	}
	// Check LHS for Koo-prefixed variable names
	for _, lhs := range s.Lhs {
		if ident, ok := lhs.(*ast.Ident); ok {
			// Match common patterns like userHandler for KooUserHandler
			if strings.Contains(strings.ToLower(ident.Name), "koo") {
				return true
			}
		}
	}

	return false
}

func removeKooRoutes(block *ast.BlockStmt) bool {
	if block == nil {
		return false
	}

	modified := false

	ast.Inspect(block, func(n ast.Node) bool {
		if compLit, ok := n.(*ast.CompositeLit); ok {
			if arrayType, ok := compLit.Type.(*ast.ArrayType); ok {
				if ident, ok := arrayType.Elt.(*ast.Ident); ok {
					if ident.Name == "route" {
						// Filter out route entries with /koo/ path
						var newElts []ast.Expr
						for _, elt := range compLit.Elts {
							if routeLit, ok := elt.(*ast.CompositeLit); ok {
								if routeHasKooPath(routeLit) {
									modified = true
									continue
								}
							}

							newElts = append(newElts, elt)
						}

						compLit.Elts = newElts
					}
				}
			}
		}

		return true
	})

	return modified
}

func routeHasKooPath(routeLit *ast.CompositeLit) bool {
	for _, elt := range routeLit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok {
				if ident.Name == "Path" {
					if lit, ok := kv.Value.(*ast.BasicLit); ok {
						if strings.Contains(lit.Value, "/koo/") {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func removeKooStructLitFields(block *ast.BlockStmt) bool {
	if block == nil {
		return false
	}

	modified := false

	ast.Inspect(block, func(n ast.Node) bool {
		if compLit, ok := n.(*ast.CompositeLit); ok {
			var newElts []ast.Expr
			for _, elt := range compLit.Elts {
				if kv, ok := elt.(*ast.KeyValueExpr); ok {
					if ident, ok := kv.Key.(*ast.Ident); ok {
						if hasKooPrefix(ident.Name) {
							modified = true
							continue
						}
					}
					// Also check if value references a Koo type/function
					if valueHasKooRef(kv.Value) {
						modified = true
						continue
					}
				}

				newElts = append(newElts, elt)
			}

			if len(newElts) != len(compLit.Elts) {
				compLit.Elts = newElts
			}
		}

		return true
	})

	return modified
}

func valueHasKooRef(expr ast.Expr) bool {
	switch v := expr.(type) {
	case *ast.Ident:
		return hasKooPrefix(v.Name)
	case *ast.CallExpr:
		if ident, ok := v.Fun.(*ast.Ident); ok {
			return hasKooPrefix(ident.Name)
		}

		if sel, ok := v.Fun.(*ast.SelectorExpr); ok {
			return hasKooPrefix(sel.Sel.Name)
		}
	case *ast.SelectorExpr:
		return hasKooPrefix(v.Sel.Name)
	}

	return false
}

func runGoimports(root string) error {
	internalPath := filepath.Join(root, "internal")
	cmd := exec.Command("goimports", "-w", internalPath) //nolint:gosec // internalPath is constructed from trusted root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
