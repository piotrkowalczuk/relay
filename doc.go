// Package relay provides IRC client implementations.
//
//	package main
//
//	import (
//		"crypto/tls"
//		"log"
//		"os"
//
//		"github.com/piotrkowalczuk/relay"
//		"github.com/piotrkowalczuk/relay/action"
//		"github.com/piotrkowalczuk/relay/freenode"
//		"github.com/sorcix/irc"
//	)
//
//	func main() {
//		conn, err := tls.Dial("tcp", freenode.Addr, &tls.Config{})
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		rel := relay.NewClientWithOpts(
//			conn,
//			&relay.User{
//				Nick:     "bot",
//				RealName: "Mr Bot",
//				Password: "password",
//				Mode:     irc.UserModeOperator,
//			},
//			&relay.ClientOpts{
//				Logger: log.New(os.Stderr, "", log.LstdFlags),
//			},
//		)
//
//		sm := relay.NewServeMux()
//		sm.Handle(irc.PRIVMSG, relay.HandleFunc(func(mw *relay.MessagesWriter, r *relay.Request){
//			mw.WriteParams(r.Receivers()...)
// 			mw.Write([]byte("Reply"))
// 		})
//
//		rel.Handle(handler)
//		rel.ListenAndReply()
//
//		for {
//			select {
//			case <-rel.Registered():
//				err := rel.Join(relay.NewChannel("#go-nuts", ""))
//				if err != nil {
//					log.Fatal(err)
//				}
//			case err := <-rel.Err():
//				log.Println("ERROR: ", err)
//			}
//		}
//	}
//
// This is a complete program that connects to FreeNode and expose one handler for all private messages.
package relay
