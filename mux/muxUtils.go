package mux

func (m *Mux) containsRoute(route string) (int, bool) {
	for i, r := range m.Routes {
		if r.URL == route {
			return i, true
		}
	}

	return -1, false
}
