package banner

import "testing"

func TestBanner(t *testing.T) {
    ban := Banner()
    if ban != "KAPOW!!!" {
        t.Errorf("Banner expected KAPOW!!!, but got %v", ban)
        t.Fail()
    }
}
