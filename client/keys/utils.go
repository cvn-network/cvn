package keys

import (
	"fmt"
	"io"

	"github.com/cosmos/cosmos-sdk/client/keys"
	cryptokeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"sigs.k8s.io/yaml"
)

// available output formats.
const (
	OutputFormatText = "text"
	OutputFormatJSON = "json"
)

type bechKeyOutFn func(k *cryptokeyring.Record) (cryptokeyring.KeyOutput, error)

func printKeyringRecord(w io.Writer, k *cryptokeyring.Record, bechKeyOut bechKeyOutFn, output string) error {
	ko, err := bechKeyOut(k)
	if err != nil {
		return err
	}

	switch output {
	case OutputFormatText:
		if err := printTextRecords(w, []cryptokeyring.KeyOutput{ko}); err != nil {
			return err
		}

	case OutputFormatJSON:
		out, err := keys.KeysCdc.MarshalJSON(ko)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintln(w, string(out)); err != nil {
			return err
		}
	}

	return nil
}

func printTextRecords(w io.Writer, kos []cryptokeyring.KeyOutput) error {
	out, err := yaml.Marshal(&kos)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, string(out)); err != nil {
		return err
	}

	return nil
}
