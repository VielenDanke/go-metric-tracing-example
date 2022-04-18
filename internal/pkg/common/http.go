package common

import "net/http"

func LiveReadyProbe(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)
}
