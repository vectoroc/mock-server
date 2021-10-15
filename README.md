mock-server
====

Http proxy that allows to mock requests by conditions and to define response policies.

Implements subset of [MockServer](https://github.com/mock-server/mockserver) API.

Implemented:

- partial support of http mocking (/expectation)
- clear / reset requests
- http & https proxy (via connect method)

TODO
====

- [x] request normalization
- [x] expectations Priority
- [x] KeyTo*Value un/marshaling
- [ ] unify logging
- [ ] benchmarks
- [ ] fulfill unit tests
- [ ] integration tests with http requests
- [ ] Body unmarshaling
- [ ] JSON match (currently it works as a text match)
- [ ] support expectation Times & TimeToLive
- [ ] support HttpRequest.KeepAlive
- [ ] support HttpResponse.ConnectionOptions
- [ ] support HttpResponse.ReasonPhrase
- [ ] models validation (github.com/asaskevich/govalidator ?)
- [ ] https ?
