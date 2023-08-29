package cmd

import (
	"context"
	"errors"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"time"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/mshindle/tidbits/leader"
)

// retryCmd represents the retry command
var leaderCmd = &cobra.Command{
	Use:   "leader",
	Short: "simple leader election using the bully algorithm",
	Long: `
In distributed systems, a leader is a concept used to manage coordination and communication among multiple
nodes or servers. In the bully algorithm, the fundamental idea is rank. It assumes that every node has a rank
within the cluster, and the leader must be the highest. So it uses the node’s rank value during the election.

There are two situations for election.

 - The system is newly initialized, so there is no leader
 - One of the nodes notices that the leader is down.
`,
	RunE: runNode,
}

func init() {
	rootCmd.AddCommand(leaderCmd)
	leaderCmd.Flags().Int("rank", 0, "node rank for runNode election")
	leaderCmd.Flags().Int("nodes", 1, "number of nodes involved in runNode election")

	_ = viper.BindPFlag("rank", leaderCmd.Flags().Lookup("rank"))
	_ = viper.BindPFlag("nodes", leaderCmd.Flags().Lookup("nodes"))
}

func runNode(cmd *cobra.Command, args []string) error {
	rank := viper.GetInt("rank")
	totalNodes := viper.GetInt("nodes")
	entry := log.WithField("rank", rank).WithField("total_nodes", totalNodes)

	if rank <= 0 || rank > totalNodes {
		entry.Error("rank is out of bounds for this node")
		return errors.New("rank should be greater than 0 and less than or equal to nodes")
	}

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	node := leader.NewNode(rank, totalNodes)
	listener, err := node.Listen(ctx)
	if err != nil {
		return err
	}
	defer func(l net.Listener) {
		_ = l.Close()
	}(listener)

	rpcServer := rpc.NewServer()
	_ = rpcServer.Register(node)
	go rpcServer.Accept(listener)

	node.ConnectToPeers()
	log.WithField("node_id", node.ID).WithField("peers", node.Peers.ToIDs()).Info("aware of peers")

	time.Sleep(5 * time.Second)
	node.Elect()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	return nil
}
