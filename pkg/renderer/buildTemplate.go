package renderer

import "html/template"

func (t *TemplateRenderer) buildTemplate(name string) (*template.Template, error) {
	filePaths, err := t.resolveTemplateFiles(name)
	if err != nil {
		return nil, err
	}

	tmpl, err := t.shared.Clone()
	if err != nil {
		return nil, err
	}

	_, err = tmpl.ParseFiles(filePaths...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
