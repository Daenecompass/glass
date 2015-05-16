package vcs

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hashicorp/errwrap"
)

var PostCheckoutTmpl = template.Must(template.New("name").Parse(`#!/bin/sh
# when checkout is a branch, start timer
if [ $3 -eq 1 ]; then
   glass start;
fi
`))

var PrepCommitTmpl = template.Must(template.New("name").Parse(`#!/bin/sh
# only add time to template and message sources
# @see http://git-scm.com/docs/githooks#_prepare_commit_msg
case "$2" in
message|template) 
	printf "$(cat $1) [$(glass split)]" > "$1" ;;
esac
`))

var PostCommitTmpl = template.Must(template.New("name").Parse(`#!/bin/sh
#always reset after commit
glass lap
`))

type Git struct {
	dir string
}

func NewGit(dir string) *Git {
	return &Git{
		dir: filepath.Join(dir, ".git"),
	}
}

func (g *Git) Name() string { return "git" }
func (g *Git) Supported() bool {
	fi, err := os.Stat(g.dir)
	if err != nil || !fi.IsDir() {
		return false
	}

	return true
}

func (g *Git) Hook() error {
	hpath := filepath.Join(g.dir, "hooks")

	//post checkout: start()
	postchf, err := os.Create(filepath.Join(hpath, "post-checkout"))
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Failed to create post-checkout '%s': {{err}}", postchf.Name()), err)
	}

	err = postchf.Chmod(0766)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Failed to make post-checkout file '%s' executable: {{err}}", hpath), err)
	}

	err = PostCheckoutTmpl.Execute(postchf, struct{}{})
	if err != nil {
		return errwrap.Wrapf("Failed to run post-checkout template: {{err}}", err)
	}

	//prepare commit msg: split()
	prepcof, err := os.Create(filepath.Join(hpath, "prepare-commit-msg"))
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Failed to create prepare-commit-msg  '%s': {{err}}", postchf.Name()), err)
	}

	err = prepcof.Chmod(0766)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Failed to make prepare-commit-msg file '%s' executable: {{err}}", hpath), err)
	}

	err = PrepCommitTmpl.Execute(prepcof, struct{}{})
	if err != nil {
		return errwrap.Wrapf("Failed to run post-commit template: {{err}}", err)
	}

	//post commit: lap()
	postcof, err := os.Create(filepath.Join(hpath, "post-commit"))
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Failed to create post-commit  '%s': {{err}}", postchf.Name()), err)
	}

	err = postcof.Chmod(0766)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Failed to make post-commit file '%s' executable: {{err}}", hpath), err)
	}

	err = PostCommitTmpl.Execute(postcof, struct{}{})
	if err != nil {
		return errwrap.Wrapf("Failed to run post-commit template: {{err}}", err)
	}

	return nil
}
