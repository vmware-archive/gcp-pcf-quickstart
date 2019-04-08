package pattern

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	cfgtypes "github.com/cloudfoundry/config-server/types"

	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/cppforlife/go-patch/patch"
)

func (t *Template) Evaluate(expectAllKeys bool) ([]byte, error) {
	template, err := t.readFile(t.Manifest)
	if err != nil {
		return []byte{}, err
	}

	tpl := boshtpl.NewTemplate(template)

	var firstToUse []boshtpl.Variables

	staticVars := boshtpl.StaticVariables{}
	ops := patch.Ops{}

	for _, file := range t.OpsFiles {
		var opDefs []patch.OpDefinition
		err = t.readYAMLFile(file, &opDefs)
		if err != nil {
			return nil, err
		}
		op, err := patch.NewOpsFromDefinitions(opDefs)
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}

	for _, file := range t.VarsFiles {
		var fileVars boshtpl.StaticVariables
		err = t.readYAMLFile(file, &fileVars)
		if err != nil {
			return nil, err
		}
		for k, v := range fileVars {
			staticVars[k] = v
		}
	}

	for k, v := range t.Vars {
		staticVars[k] = v
	}

	firstToUse = append(firstToUse, staticVars)

	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)
	store := &boshcmd.VarsFSStore{FS: fs}

	if t.VarsStore != "" {
		err := store.UnmarshalFlag(t.VarsStore)
		if err != nil {
			return []byte{}, err
		}
	}

	if store.IsSet() {
		firstToUse = append(firstToUse, store)
	}

	vars := boshtpl.NewMultiVars(firstToUse)

	if store.IsSet() {
		store.ValueGeneratorFactory = cfgtypes.NewValueGeneratorConcrete(boshcmd.NewVarsCertLoader(vars))
	}

	evalOpts := boshtpl.EvaluateOpts{
		UnescapedMultiline: true,
		ExpectAllKeys:      expectAllKeys,
		ExpectAllVarsUsed:  false,
	}

	bytes, err := tpl.Evaluate(vars, ops, evalOpts)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (t *Template) readFile(file string) ([]byte, error) {
	if filepath.Ext(file) == "" {
		file = fmt.Sprintf("%s.yml", file)
	}
	f, err := t.Store.Open(file)
	if err != nil {
		return []byte{}, err
	}
	return ioutil.ReadAll(f)
}

func (t *Template) readYAMLFile(file string, dataType interface{}) error {
	payload, err := t.readFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(payload, dataType)
	if err != nil {
		return err
	}
	return nil
}
