package codegen

import (
	"bytes"
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/grafana/pkg/kindsys"
	"github.com/grafana/thema"
)

type OneToOne codejen.OneToOne[*DeclForGen]
type OneToMany codejen.OneToMany[*DeclForGen]
type ManyToOne codejen.ManyToOne[*DeclForGen]
type ManyToMany codejen.ManyToMany[*DeclForGen]

// ForGen is a codejen input transformer that converts a pure kindsys.SomeDecl into
// a DeclForGen by binding its contained lineage.
func ForGen(rt *thema.Runtime, decl *kindsys.SomeDecl) (*DeclForGen, error) {
	lin, err := decl.BindKindLineage(rt)
	if err != nil {
		return nil, err
	}

	return &DeclForGen{
		SomeDecl: decl,
		lin:      lin,
	}, nil
}

// DeclForGen wraps [kindsys.SomeDecl] to provide trivial caching of
// the lineage declared by the kind (nil for raw kinds).
type DeclForGen struct {
	*kindsys.SomeDecl
	lin thema.Lineage
}

func (decl *DeclForGen) Lineage() thema.Lineage {
	return decl.lin
}

func SlashHeaderMapper(maingen string) codejen.FileMapper {
	return func(f codejen.File) (codejen.File, error) {
		b := new(bytes.Buffer)
		fmt.Fprintf(b, headerTmpl, maingen, f.FromString())
		fmt.Fprint(b, string(f.Data))
		f.Data = b.Bytes()
		return f, nil
	}
}

var headerTmpl = `// THIS FILE IS GENERATED. EDITING IS FUTILE.
//
// Generated by:
//     %s
// Using jennies:
//     %s
//
// Run 'make gen-cue' from repository root to regenerate.

`
