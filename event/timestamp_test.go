package event

import (
	"testing"
	"time"
)

func TestTimeStamp_UnmarshalJSON(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		timeStamp := &TimeStamp{}
		err := timeStamp.UnmarshalJSON([]byte("1355517536.000001"))
		if err != nil {
			t.Fatalf("Unexpected error is returned: %s", err.Error())
		}

		if timeStamp.OriginalValue != "1355517536.000001" {
			t.Errorf("Expected original value is not returned: %s", timeStamp.OriginalValue)
		}

		if timeStamp.Time != time.Unix(1355517536, 0) {
			t.Errorf("Expected timestamp is not set: %s", timeStamp.Time)
		}
	})

	t.Run("invalid string", func(t *testing.T) {
		timeStamp := &TimeStamp{}
		err := timeStamp.UnmarshalJSON([]byte("abc"))
		if err == nil {
			t.Fatal("Expected error is not returned.")
		}
	})
}

func TestTimeStamp_String(t *testing.T) {
	origValue := "1355517536.000001"
	timeStamp := &TimeStamp{
		OriginalValue: origValue,
	}

	if timeStamp.String() != origValue {
		t.Errorf("Expected original value is not returned: %s", timeStamp.OriginalValue)
	}
}

func TestTimeStamp_MarshalText(t *testing.T) {
	origValue := "1355517536.000001"
	timeStamp := &TimeStamp{
		OriginalValue: origValue,
	}

	b, err := timeStamp.MarshalText()
	if err != nil {
		t.Fatalf("Unexpected error is returned: %s", err.Error())
	}

	if string(b) != origValue {
		t.Errorf("Expected text is not retuend: %s", b)
	}
}
