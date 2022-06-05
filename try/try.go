package try

import "log"

func Try[K any](v K, err error) K {
	if err != nil {
		log.Fatal(err)
	}
	return v
}
