package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PrometheusConfig struct {
	Global struct {
		ScrapeInterval     string `yaml:"scrape_interval"`
		ScrapeTimeout      string `yaml:"scrape_timeout"`
		EvaluationInterval string `yaml:"evaluation_interval"`
	} `yaml:"global"`
	ScrapeConfigs []ScrapeConfigs `yaml:"scrape_configs"`
}
type ScrapeConfigs struct {
	JobName        string        `yaml:"job_name"`
	StaticConfigs  StaticConfigs `yaml:"static_configs"`
	MetricsPath    string        `yaml:"metrics_path,omitempty"`
	BasicAuth      BasicAuth     `yaml:"basic_auth,omitempty"`
	ScrapeInterval string        `yaml:"scrape_interval,omitempty"`
	ScrapeTimeout  string        `yaml:"scrape_timeout,omitempty"`
}
type StaticConfigs struct {
	Targets []string `yaml:"targets"`
}

type BasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Jobs struct {
	Jobs     []JsonData `json:"jobs"`
	Password string     `json:"password"`
	Target   string     `json:"target"`
	Interval string     `json:"interval"`
	Timeout  string     `json:"timeout"`
}
type JsonData struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Interval string `json:"interval"`
}

func (j *Jobs) ConvertToYml() error {
	var p *PrometheusConfig
	p.Global.ScrapeInterval = j.Interval
	p.Global.ScrapeTimeout = j.Timeout
	p.Global.EvaluationInterval = j.Interval

	for _, job := range j.Jobs {
		scrapeConfig := ScrapeConfigs{
			JobName:        job.Name,
			MetricsPath:    job.Path,
			StaticConfigs:  StaticConfigs{[]string{j.Target}},
			BasicAuth:      BasicAuth{"checker", j.Password},
			ScrapeInterval: job.Interval,
			ScrapeTimeout:  j.Timeout,
		}
		p.ScrapeConfigs = append(p.ScrapeConfigs, scrapeConfig)
	}
	p.GenerateConfig("")
	return nil
}

func (p *PrometheusConfig) GenerateConfig(filename string) error {
	fmt.Println(filename, p)
	return nil
}

func JsonParser(w http.ResponseWriter, r *http.Request) {
	var j Jobs
	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
