/*

MIT License

Copyright (c) 2017 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package udp

import (
	"fmt"
	"net"

	"github.com/piot/flux-go/src/endpoint"
)

type Communication struct {
	hostAddr      *net.UDPAddr
	udpConnection *net.UDPConn
}

func (s *Communication) WriteToUDP(addr *endpoint.Endpoint, octets []byte) error {
	a := addr.UDPAddr()
	sentOctets, writeErr := s.udpConnection.WriteToUDP(octets, a)
	if writeErr != nil {
		fmt.Printf("UDP Write failed '%v' %v %v\n", writeErr, a, octets)
		return writeErr
	}
	if sentOctets != len(octets) {
		sentOctetsErr := fmt.Errorf("didn't send all octets:%v expected:%v", sentOctets, len(octets))
		fmt.Printf("UDP Write failed '%v' %v %v\n", sentOctetsErr, a, octets)
		return sentOctetsErr
	}
	return nil
}

func (s *Communication) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	n, addr, err := s.udpConnection.ReadFromUDP(b)
	return n, addr, err
}

func (s *Communication) HostAddr() *net.UDPAddr {
	return s.hostAddr
}

func NewServerCommunication(listenPort int) (*Communication, error) {
	portString := fmt.Sprintf(":%d", listenPort)
	localAddr, localAddrErr := net.ResolveUDPAddr("udp", portString)
	if localAddrErr != nil {
		return nil, localAddrErr
	}
	serverConnection, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}

	fmt.Printf("Listening to %s\n", portString)
	return &Communication{udpConnection: serverConnection}, nil
}

func NewClientCommunication(host string) (*Communication, error) {
	fmt.Printf("Connecting to %v", host)
	serverAddr, serverAddrErr := net.ResolveUDPAddr("udp", host)
	if serverAddrErr != nil {
		return nil, serverAddrErr
	}
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, err
	}
	return &Communication{udpConnection: conn, hostAddr: serverAddr}, nil
}
