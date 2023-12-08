package repository

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"os/exec"

	"actlabs-auth/entity"
)

type labRepository struct{}

func NewLabRepository() entity.LabRepository {
	return &labRepository{}
}

func (l *labRepository) GetEnumerationResults(typeOfLab string, includeVersions bool) (entity.EnumerationResults, error) {
	er := entity.EnumerationResults{}

	// URL of the container to list the blobs
	url := "https://" + entity.StorageAccountName + ".blob.core.windows.net/repro-project-" + typeOfLab + "" + entity.SasToken + "&restype=container&comp=list"

	if includeVersions {
		url = url + "&include=versions"
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return er, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	err = xml.Unmarshal(body, &er)
	if err != nil {
		return er, err
	}

	return er, nil
}

func (l *labRepository) GetLab(typeOfLab string, labId string, versionId string) (entity.LabType, error) {
	lab := entity.LabType{}

	// if labId is doesn't end in .json then append it
	if len(labId) < 5 || labId[len(labId)-5:] != ".json" {
		labId = labId + ".json"
	}

	// URL of the blob with SAS token.
	url := "https://" + entity.StorageAccountName + ".blob.core.windows.net/repro-project-" + typeOfLab + "/" + labId + "" + entity.SasToken

	if versionId != "" {
		url = url + "&versionid=" + versionId
	}

	resp, err := http.Get(url)
	if err != nil {
		return lab, err
	}

	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(bodyBytes, &lab); err != nil {
		return lab, err
	}

	return lab, nil
}

// Create or add a new version of lab.
func (l *labRepository) UpsertLab(labId string, lab string, typeOfLab string) error {
	_, err := exec.Command("bash", "-c", "echo '"+lab+"' | az storage blob upload --data @- -c repro-project-"+typeOfLab+" -n "+labId+".json --account-name "+entity.StorageAccountName+" --sas-token '"+entity.SasToken+"' --overwrite").Output()
	return err
}

// Deletes all versions of lab.
func (l *labRepository) DeleteLab(typeOfLab string, labId string) error {
	_, err := exec.Command("bash", "-c", "az storage blob delete -c repro-project-"+typeOfLab+" -n "+labId+".json --account-name "+entity.StorageAccountName+" --sas-token '"+entity.SasToken+"'").Output()
	return err
}
