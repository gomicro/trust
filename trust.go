package trust

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

// Pool is a representation of a CA Cert Pool that can have multiple items
// appended into it
type Pool struct {
	pool  *x509.CertPool
	files []string
}

// New returns a newly initialized CA Pool that can then have additional actions
// performed on it
func New() *Pool {
	return &Pool{
		pool: x509.NewCertPool(),
	}
}

// AddCAFile adds the specified files to the list of CA files which will be
// appended to the resulting pool
func (p *Pool) AddCAFile(files ...string) {
	p.files = append(p.files, files...)
}

// CACerts builds an X.509 certificate pool containing the Mozilla CA
// Certificate bundle. Returns nil on error along with an appropriate error
// code.
func (p *Pool) CACerts() (*x509.CertPool, error) {
	if p.pool == nil {
		p.pool = x509.NewCertPool()
	}

	err := p.appendDefaultCerts()
	if err != nil {
		return nil, fmt.Errorf("failed to append default certs to pool: %v", err.Error())
	}

	err = p.appendFileCerts()
	if err != nil {
		return nil, fmt.Errorf("failed to append file certs to pool: %v", err.Error())
	}

	return p.pool, nil
}

func (p *Pool) appendDefaultCerts() error {
	ok := p.pool.AppendCertsFromPEM([]byte(globalPemCerts))
	if !ok {
		return fmt.Errorf("failed to append global CAs to cert pool")
	}

	return nil
}

func (p *Pool) appendFileCerts() error {
	for _, file := range p.files {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read cert file: %v", err.Error())
		}

		ok := p.pool.AppendCertsFromPEM(b)
		if !ok {
			return fmt.Errorf("failed to append cert file (%v) to cert pool", file)
		}
	}

	return nil
}
