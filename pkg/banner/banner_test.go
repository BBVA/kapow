package banner

import "testing"

func TestBanner(t *testing.T) {
    ban := Banner("0.0.0")
    if ban == "" {
        t.Errorf("Banner expected KAPOW!!!, but got %v", ban)
        t.Fail()
    }
}
