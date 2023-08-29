package leader

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"time"

	"github.com/apex/log"
)

type Node struct {
	ID         string
	Addr       string
	Rank       int
	TotalNodes int
	Peers      *Peers
	eventBus   *EventBus
}

const (
	basePort     = 6000
	maxRetries   = 10
	waitInMS     = 250
	pingInterval = 3 // how often to ping leader is still active...
)

// lookupNode is a convenience function to determine node information that
// service discovery would have provided for us...
func lookupNode(rank int) (string, string) {
	id := fmt.Sprintf("node-%02d", rank)
	addr := fmt.Sprintf("%s:%d", id, basePort+rank)

	return id, addr
}

func NewNode(rank int, totalNodes int) *Node {
	id, addr := lookupNode(rank)
	node := &Node{
		ID:         id,
		Addr:       addr,
		Rank:       rank,
		TotalNodes: totalNodes,
		Peers:      NewPeers(),
		eventBus:   NewEventBus(),
	}
	node.eventBus.Subscribe(LeaderElected, node.PingLeader)

	return node
}

func (n *Node) Listen(ctx context.Context) (net.Listener, error) {
	var lc net.ListenConfig
	return lc.Listen(ctx, "tcp", n.Addr)
}

func (n *Node) ConnectToPeers() {
	for r := 1; r <= n.TotalNodes; r++ {
		peerID, peerAddr := lookupNode(r)
		if n.IsItself(peerID) {
			continue
		}
		entry := log.WithField("peer", peerID)

		client, err := n.connect(peerAddr)
		if err != nil {
			entry.WithError(err).Error("failed to connect to peer... skipping")
			continue
		}

		pingMessage := Message{FromPeerID: n.ID, Type: PING}
		reply, err := n.CommunicateWithPeer(client, pingMessage)
		if err != nil {
			entry.WithError(err).Error("failed handle message call... skipping")
			continue
		}

		if reply.IsPongMessage() {
			entry.Info("got pong message from peer")
			n.Peers.Add(reply.FromPeerID, reply.Rank, client)
		}
	}
}

func (n *Node) connect(peerAddr string) (client *rpc.Client, err error) {
	for i := 0; i < maxRetries; i++ {
		client, err = rpc.Dial("tcp", peerAddr)
		if err == nil {
			return
		}

		log.WithError(err).WithField("addr", peerAddr).WithField("count", i).Info("error dialing rpc")
		time.Sleep(waitInMS * time.Millisecond)
	}
	return
}

func (n *Node) CommunicateWithPeer(client *rpc.Client, args Message) (Message, error) {
	var reply Message

	err := client.Call("Node.HandleMessage", args, &reply)
	return reply, err
}

func (n *Node) HandleMessage(args Message, reply *Message) error {
	reply.FromPeerID = n.ID
	reply.Rank = n.Rank

	switch args.Type {
	case ELECTION:
		reply.Type = ALIVE
	case ELECTED:
		leaderID := args.FromPeerID
		log.WithField("leader_id", leaderID).Info("election is done - new leader set")
		n.eventBus.Emit(LeaderElected, leaderID)
		reply.Type = OK
	case PING:
		reply.Type = PONG
	}

	return nil
}

func (n *Node) Elect() {
	isHighestRankedNodeAvailable := false

	for _, peer := range n.Peers.ToList() {
		if n.IsRankHigher(peer.Rank) {
			continue
		}
		entry := log.WithField("peer", peer.ID).WithField("node", n.ID)

		entry.Info("send ELECTION message to peer")
		electionMessage := Message{FromPeerID: n.ID, Rank: n.Rank, Type: ELECTION}
		reply, err := n.CommunicateWithPeer(peer.Client, electionMessage)
		if err != nil {
			entry.WithError(err).Error("failed handle message call... skipping")
			continue
		}

		if reply.IsAliveMessage() {
			isHighestRankedNodeAvailable = true
		}
	}

	if !isHighestRankedNodeAvailable {
		electedMessage := n.createMessage(ELECTED)
		n.BroadcastMessage(electedMessage)
		log.WithField("node", n.ID).Info("node is a new leader")
	}
}

func (n *Node) BroadcastMessage(args Message) {
	for _, peer := range n.Peers.ToList() {
		_, _ = n.CommunicateWithPeer(peer.Client, args)
	}
}

func (n *Node) PingLeader(event Event, payload any) {
	if event != LeaderElected {
		log.WithField("event", event).Error("invalid event for PingLeader")
		return
	}

	leaderID := payload.(string)
	entry := log.WithField("leader_id", leaderID)
	leader := n.Peers.Get(leaderID)
	if leader == nil {
		entry.Error("no peer with specified leader_id")
		return
	}

	for {
		pingMessage := n.createMessage(PING)
		reply, err := n.CommunicateWithPeer(leader.Client, pingMessage)
		if err != nil {
			entry.Info("leader is down, new election about to start!")
			n.Peers.Delete(leaderID)
			n.Elect()
			return
		}
		if reply.IsPongMessage() {
			entry.Info("leader sent pong message")
			time.Sleep(pingInterval * time.Second)
		}
	}
}

func (n *Node) IsItself(id string) bool {
	return n.ID == id
}

func (n *Node) IsRankHigher(rank int) bool {
	return n.Rank > rank
}

func (n *Node) createMessage(mt MessageType) Message {
	return Message{
		FromPeerID: n.ID,
		Rank:       n.Rank,
		Type:       mt,
	}
}
