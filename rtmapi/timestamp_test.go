package rtmapi

import (
	"bytes"
	"github.com/oklahomer/golack/v2/event"
	"strconv"
	"testing"
	"time"
)

func TestTimeStamp_UnmarshalJSON(t *testing.T) {
	timeStamp := &event.TimeStamp{}
	if err := timeStamp.UnmarshalJSON([]byte("1355517536.000001")); err != nil {
		t.Fatalf("error on unmarshal slack timestamp: %s.", err.Error())
	}

	expectedTime := time.Unix(1355517536, 0)

	if !timeStamp.Time.Equal(expectedTime) {
		t.Errorf("unmarshaled time is wrong: %s. expected: %s.", timeStamp.Time.String(), expectedTime.String())
	}
}

func TestTimeStamp_UnmarshalJSON_Number(t *testing.T) {
	timeStamp := &event.TimeStamp{}
	if err := timeStamp.UnmarshalJSON([]byte(strconv.Itoa(1355517536))); err != nil {
		t.Fatalf("error on unmarshal slack timestamp: %s.", err.Error())
	}

	expectedTime := time.Unix(1355517536, 0)

	if !timeStamp.Time.Equal(expectedTime) {
		t.Errorf("unmarshaled time is wrong: %s. expected: %s.", timeStamp.Time.String(), expectedTime.String())
	}
}

func TestTimeStamp_UnmarshalText_Malformed(t *testing.T) {
	invalidInput := "FooBar"
	timeStamp := &event.TimeStamp{}
	if err := timeStamp.UnmarshalJSON([]byte(invalidInput)); err == nil {
		t.Error("error is not returned.")
	}
}

func TestTimeStamp_MarshalText(t *testing.T) {
	now := time.Now()
	slackTimeStamp := strconv.Itoa(now.Second()) + ".123"
	timeStamp := &event.TimeStamp{Time: now, OriginalValue: slackTimeStamp}

	b, err := timeStamp.MarshalText()

	if err != nil {
		t.Fatalf("unexpected error on timestamp marshal: %s.", err.Error())
	}

	if !bytes.Equal(b, []byte(slackTimeStamp)) {
		t.Errorf("marshaled value is wrong: %s. expected: %s.", b, slackTimeStamp)
	}
}
