package mqtt

import (
	"bytes"
	"html/template"
	"regexp"
	"strings"
)

func (t *Topic) MatchTopics(topic string) bool {
	s := strings.ReplaceAll(t.Name, "+", "[^/]+")
	return regexp.MustCompile(s).MatchString(topic)
}

func MqttWildcard(topic string, wildcard string) []string {
	if topic == wildcard {
		return []string{}
	} else if wildcard == "#" {
		return []string{topic}
	}
	res := []string{}
	t := strings.Split(topic, "/")
	w := strings.Split(wildcard, "/")
	i := 0
	for lt := len(t); i < lt; i++ {
		if w[i] == "+" {
			res = append(res, t[i])
		} else if w[i] == "#" {
			res = append(res, strings.Join(t[i:], "/"))
			return res
		} else if w[i] != t[i] {
			return nil
		}
	}
	if w[i] == "#" {
		i += 1
	}
	if i == len(w) {
		return res
	}
	return nil
}

func TopicWithTemplate(topic string, mapping map[string]string) (string, error) {
	t := template.Must(template.New("").Parse(topic))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, mapping); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
