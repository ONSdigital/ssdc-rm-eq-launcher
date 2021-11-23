package main // import "github.com/ONSdigital/ssdc-rm-eq-launcher"

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/ssdc-rm-eq-launcher/settings"
	"io/ioutil"
	"strings"

	"html"
	"log"
	"net/http"

	"github.com/ONSdigital/ssdc-rm-eq-launcher/authentication"
	"github.com/gofrs/uuid"
)

type UacResponse struct {
	Active                  bool   `json:"active"`
	CollectionInstrumentUrl string `json:"collectionInstrumentUrl"`
	Qid                     string `json:"qid"`
	CaseId                  string `json:"caseId"`
}

func getAccountServiceURL(r *http.Request) string {
	forwardedProtocol := r.Header.Get("X-Forwarded-Proto")

	requestProtocol := "http"

	if forwardedProtocol != "" {
		requestProtocol = forwardedProtocol
	}

	return fmt.Sprintf("%s://%s",
		requestProtocol,
		html.EscapeString(r.Host))
}

func checkUac(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	uac := strings.ToUpper(r.FormValue("uac"))

	handleUac(w, r, uac)
}

func launchRequest(w http.ResponseWriter, r *http.Request) {
	urlValues := r.URL.Query()
	uac := strings.ToUpper(urlValues.Get("uac"))

	handleUac(w, r, uac)
}

func handleUac(w http.ResponseWriter, r *http.Request, uac string) {
	hash := sha256.New()
	hash.Write([]byte(uac))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	caseAPIURL := settings.Get("CASE_API_URL")

	response, err := http.Get(caseAPIURL + "/uacs/" + mdStr)
	if err != nil || response.StatusCode == 404 {
		http.Redirect(w, r, "baduac.html", 302)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var uacResponse UacResponse
	json.Unmarshal(responseData, &uacResponse)

	if uacResponse.Active {
		launchEq(w, r, uacResponse.CollectionInstrumentUrl, uacResponse.Qid)
	} else {
		http.Redirect(w, r, "inactiveuac.html", 302)
	}
}

func launchEq(w http.ResponseWriter, r *http.Request, collectionInstrumentUrl string, qid string) {
	hostURL := settings.Get("SURVEY_RUNNER_URL")
	accountServiceURL := getAccountServiceURL(r)
	AccountServiceLogOutURL := getAccountServiceURL(r)
	urlValues := r.URL.Query()
	defaultValues := authentication.GetDefaultValues()
	log.Println("Launch request received", collectionInstrumentUrl)

	urlValues.Add("ru_ref", defaultValues["ru_ref"])
	collectionExerciseSid, _ := uuid.NewV4()
	caseID, _ := uuid.NewV4()
	urlValues.Add("collection_exercise_sid", collectionExerciseSid.String())
	urlValues.Add("case_id", caseID.String())
	urlValues.Add("response_id", qid)
	urlValues.Add("language_code", defaultValues["language_code"])

	token, err := authentication.GenerateTokenFromDefaults(collectionInstrumentUrl, accountServiceURL, AccountServiceLogOutURL, urlValues)
	if err != "" {
		http.Error(w, err, 400)
		return
	}

	http.Redirect(w, r, hostURL+"/session?token="+token, 302)
}

func main() {
	http.HandleFunc("/check-uac", checkUac)
	http.HandleFunc("/launch", launchRequest)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.ListenAndServe(":8000", nil)
}
