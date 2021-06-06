package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kyokomi/emoji"
)

const metadataURL = "http://metadata.google.internal/computeMetadata/v1/"

var emojilist []string

// Location is the data about this app
type Location struct {
	ProjectID         string      `json:"project_id"`
	Zone              string      `json:"zone"`
	NodeName          string      `json:"node_name"`
	ClusterName       string      `json:"cluster_name"`
	HostHeader        string      `json:"host_header"`
	PodName           string      `json:"pod_name"`
	PodNameEmoji      string      `json:"pod_name_emoji"`
	Timestamp         time.Time   `json:"timestamp"`
	PodNamespace      string      `json:"pod_namespace"`
	PodIP             string      `json:"pod_ip"`
	PodServiceAccount string      `json:"pod_service_account"`
	Metadata          string      `json:"metadata"`
	BackendService    interface{} `json:"backend_result,omitempty"`
}

func main() {
	log.Println("whereami staring")

	emojilist = createEmojiList()

	r := mux.NewRouter()
	r.HandleFunc("/", EnvironmentHandler)
	r.HandleFunc("/healthz", HealthHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// EnvironmentHandler returns info about the environment
func EnvironmentHandler(w http.ResponseWriter, r *http.Request) {

	env := Location{}

	// GCP Project ID
	env.ProjectID = getMetadata("project/project-id")
	// GCP Zone
	env.Zone = getMetadata("instance/zone")
	// Node Name
	env.NodeName = getMetadata("instance/hostname")
	// Cluster name
	env.ClusterName = getMetadata("instance/attributes/cluster-name")
	// Host header
	env.HostHeader = r.Host
	// Pod Name
	hostname, _ := os.Hostname()
	env.PodName = hostname
	// Pod Name Emoji
	h := fnv.New32a()
	h.Write([]byte(hostname))
	//log.Printf("%v %% %v = %v", h.Sum32(), len(emojilist), h.Sum32()%uint32(len(emojilist)))
	env.PodNameEmoji = emojilist[h.Sum32()%uint32(len(emojilist))]
	// Timestamp
	env.Timestamp = time.Now()
	// Pod Namespace
	env.PodNamespace = os.Getenv("POD_NAMESPACE")
	// Pod IP
	env.PodIP = os.Getenv("POD_IP")
	// Pod Service Account
	env.PodServiceAccount = os.Getenv("POD_SERVICE_ACCOUNT")
	// Metadata
	env.Metadata = os.Getenv("METADATA")
	// Backend
	var backend bool
	backend, _ = strconv.ParseBool(os.Getenv("BACKEND_ENABLED"))
	if backend {
		log.Println("TBD would call backend here")
	}

	json, err := json.Marshal(&env)
	if err != nil {
		log.Println("Unable to marshal Location struct")
		http.Error(w, "unable to construct location", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(json)
}

// HealthHandler is a health endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

// createEmojiList returns a slice of emoji
func createEmojiList() []string {
	var emojilist []string

	/*
		// create emoji list
		emojii := [][]int{
			// Emoticons icons.
			{128513, 128591},
			// Dingbats.
			{9986, 10160},
			// Transport and map symbols.
			{128640, 128704},
		}
		for _, v := range emojii {
			for x := v[0]; x < v[1]; x++ {
				str := html.UnescapeString("&#" + strconv.Itoa(x) + ";")
				emojilist = append(emojilist, str)
			}
		}
	*/

	// create emoji list
	emojii := emoji.CodeMap()
	//emojilist = make([]string, len(emojii))
	for _, v := range emojii {
		if v != "" {
			emojilist = append(emojilist, v)
		}
	}

	return emojilist
}

// getMetadata retrieves the GCP Compute Metadata for the given path
func getMetadata(metadataPath string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", metadataURL, metadataPath), nil)
	req.Header.Add("Metadata-Flavor", "Google")
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	d, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return fmt.Sprintf("%s", d)
}
