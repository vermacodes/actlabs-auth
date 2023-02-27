package repository

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"os/exec"

	"actlabs-auth/entity"

	"golang.org/x/exp/slog"
)

type labRepository struct{}

func NewLabRepository() entity.LabRepository {
	return &labRepository{}
}

func (l *labRepository) GetEnumerationResults(typeOfLab string) (entity.EnumerationResults, error) {
	er := entity.EnumerationResults{}

	// URL of the container to list the blobs
	url := "https://" + entity.StorageAccountName + ".blob.core.windows.net/repro-project-" + typeOfLab + "" + entity.SasToken + "&restype=container&comp=list"

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

func (l *labRepository) GetLab(name string, typeOfLab string) (entity.LabType, error) {
	lab := entity.LabType{}

	// URL of the blob with SAS token.
	url := "https://" + entity.StorageAccountName + ".blob.core.windows.net/repro-project-" + typeOfLab + "/" + name + "" + entity.SasToken

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

func (l *labRepository) AddLab(labId string, lab string, typeOfLab string) error {
	_, err := exec.Command("bash", "-c", "echo '"+lab+"' | az storage blob upload --data @- -c repro-project-"+typeOfLab+"s -n "+labId+".json --account-name "+entity.StorageAccountName+" --sas-token '"+entity.SasToken+"' --overwrite").Output()
	return err
}

func (l *labRepository) DeleteLab(labId string, typeOfLab string) error {
	slog.Info("Command" + "az storage blob delete -c repro-project-" + typeOfLab + "s -n " + labId + ".json --account-name " + entity.StorageAccountName + " --sas-token '" + entity.SasToken + "'")
	_, err := exec.Command("bash", "-c", "az storage blob delete -c repro-project-"+typeOfLab+"s -n "+labId+".json --account-name "+entity.StorageAccountName+" --sas-token '"+entity.SasToken+"'").Output()
	return err
}
