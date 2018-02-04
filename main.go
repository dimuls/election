package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/howeyc/gopass"
	"github.com/someanon/election/contract"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Usage = "Election smart contract controller"
	app.Version = "0.1.0"
	app.Author = "Vadim Chernov"
	app.Email = "dimuls@yandex.ru"

	app.Commands = []cli.Command{
		{
			Name:  "deploy",
			Usage: "deploy election smart contract to the ethereum block chain",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ipc",
					Usage: "ethereum client's IPC path",
				},
				cli.StringFlag{
					Name:  "key",
					Usage: "chairman's encrypted key path",
				},
			},
			Action: deploy,
		},
		{
			Name:  "add-voters",
			Usage: "add voters to the election smart contract",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ipc",
					Usage: "ethereum client's IPC path",
				},
				cli.StringFlag{
					Name:  "key",
					Usage: "chairman's encrypted key path",
				},
				cli.StringFlag{
					Name:  "contract-addr",
					Usage: "election smart contract's address",
				},
				cli.StringFlag{
					Name:  "voter-addrs",
					Usage: "comma separated voters' addresses to add",
				},
			},
			Action: addVoters,
		},
		{
			Name:  "web-server",
			Usage: "run web server with ui to vote in the election smart contract",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "bind-addr",
					Usage: "web server address to bind",
					Value: "127.0.0.1:8080",
				},
				cli.StringFlag{
					Name:  "contract-addr",
					Usage: "election smart contract's address",
				},
			},
			Action: webServer,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}

func deploy(c *cli.Context) error {
	conn, err := ethclient.Dial(c.String("ipc"))
	if err != nil {
		log.Fatalf("Failed to connect to the ethereum client: %v", err)
	}
	keyFile, err := os.Open(c.String("key"))
	if err != nil {
		log.Fatalf("Failed to open key file: %v", err)
	}
	fmt.Print("Enter passphrase: ")
	passBytes, err := gopass.GetPasswd()
	if err != nil {
		log.Fatalf("Failed to get passphrase: %v", err)
	}
	auth, err := bind.NewTransactor(keyFile, string(passBytes))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	address, tx, _, err := contract.DeployElection(auth, conn)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	log.Printf("Contract pending deploy: 0x%x\n", address)
	log.Printf("Transaction waiting to be mined: 0x%x\n", tx.Hash())
	for {
		time.Sleep(1 * time.Second)
		_, isPending, err := conn.TransactionByHash(context.TODO(), tx.Hash())
		if err != nil {
			log.Printf("Failed to get transaction by hash: %v\n", err)
			continue
		}
		if !isPending {
			log.Println("Transaction mined, contract deployed.")
			break
		}
	}
	return nil
}

func addVoters(c *cli.Context) error {
	conn, err := ethclient.Dial(c.String("ipc"))
	if err != nil {
		log.Fatalf("Failed to connect to the ethereum client: %v", err)
	}
	keyFile, err := os.Open(c.String("key"))
	if err != nil {
		log.Fatalf("Failed to open key file: %v", err)
	}
	fmt.Print("Enter passphrase: ")
	passBytes, err := gopass.GetPasswd()
	if err != nil {
		log.Fatalf("Failed to get passphrase: %v", err)
	}
	auth, err := bind.NewTransactor(keyFile, string(passBytes))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	el, err := contract.NewElection(common.HexToAddress(c.String("contract-addr")), conn)
	if err != nil {
		log.Fatalf("Failed to create election contract: %v", err)
	}
	var as []common.Address
	for _, a := range strings.Split(c.String("voter-addrs"), ",") {
		as = append(as, common.HexToAddress(a))
	}
	tx, err := el.AddVoters(auth, as)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	log.Printf("Transaction waiting to be mined: 0x%x\n", tx.Hash())
	for {
		time.Sleep(1 * time.Second)
		_, isPending, err := conn.TransactionByHash(context.TODO(), tx.Hash())
		if err != nil {
			log.Printf("Failed to get transaction by hash: %v\n", err)
			continue
		}
		if !isPending {
			log.Println("Transaction mined, voter added.")
			break
		}
	}
	return nil
}

func webServer(c *cli.Context) error {
	var t *template.Template
	t = template.New("index.html")
	if _, err := t.Parse(indexPage); err != nil {
		log.Fatalf("Failed to parse indexPage template: %v", err)
	}
	data := struct {
		ContractABI  string
		ContractAddr string
	}{
		ContractABI:  contract.ElectionABI,
		ContractAddr: c.String("contract-addr"),
	}
	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, data); err != nil {
		log.Fatalf("Failed to execute index.html template: %v", err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s", r.RemoteAddr, r.Method, r.RequestURI)
		w.Write(buf.Bytes())
	})
	if err := http.ListenAndServe(c.String("bind-addr"), nil); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
	return nil
}

const indexPage = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Выборы 2018</title>
    <style>
        .active {
            color: blue;
        }
        .error {
            color: red;
        }
        .vote {
            display: block;
            margin: 5px 0;
        }
        .success {
            color: green;
        }
    </style>
</head>
<body>
    <h1>Выборы 2018</h1>
    <div data-bind="visible: false" class="active">
        Загрузка...
    </div>
    <div data-bind="visible: !web3" style="display: none">
        Используйте <a href="https://metamask.io/">Metamask</a> или <a href="https://ethereum.org/">Mist</a> для доступа к выборам
    </div>
    <div data-bind="if: web3, visible: web3" style="display: none">
        <div data-bind="foreach: errors" class="error">
            <div data-bind="text: $data"></div>
        </div>
        <div data-bind="ifnot: errors().length">
            <div data-bind="if: checkingVoteAbility" class="active">
                Проверка возможности голосовать...
            </div>
            <div data-bind="if: !checkingVoteAbility() && !ableToVote()">
                <div data-bind="ifnot: voter">
                    Вы не имеет права голосовать
                </div>
                <div data-bind="if: voted">
                    Вы уже голосовали
                </div>
            </div>
            <div data-bind="if: ableToVote">
                <div data-bind="ifnot: candidatesLoaded" class="active">
                    Загрузка кандидатов...
                </div>
                <div data-bind="if: candidatesLoaded">
                    <div data-bind="ifnot: candidateVoted">
                        <div data-bind="foreach: candidates">
                            <button data-bind="text: name, click: vote" class="vote"></button>
                        </div>
                    </div>
                    <div data-bind="with: candidateVoted" class="success">
                        Вы успешно проголосовали за <span data-bind="text: name"></span>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/lodash.js/4.17.4/lodash.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/knockout/3.4.2/knockout-min.js"></script>
    <script>window.onload = function() {

window.app = {{.}};

if (typeof web3 === 'undefined') {
    return ko.applyBindings(app);
}

app.web3 = new Web3(web3.currentProvider);
app.ct = app.web3.eth.contract(JSON.parse(app.ContractABI)).at(app.ContractAddr);
app.errors = ko.observableArray();
app.voter = ko.observable(null);
app.voted = ko.observable(null);
app.checkingVoteAbility = ko.pureComputed(function() {
    var voter = app.voter();
    var voted = app.voted();
    switch (voter) {
        case null:
            return true;
        case false:
            return false;
        case true:
            return voted === null;
    }
});
app.ableToVote = ko.pureComputed(function() {
    var checked = !app.checkingVoteAbility();
    var voter = app.voter();
    var voted = app.voted();
    return checked && voter && !voted
});
app.candidates = ko.observableArray();
app.candidatesCount = ko.observable(null);
app.candidatesLoaded = ko.pureComputed(function() {
    var c = app.candidatesCount();
    var cs = app.candidates();
    return c !== null && c === cs.length
});
app.candidatesLoaded.subscribe(function(loaded) {
    if (loaded) {
        console.log("all candidates loaded");
        app.candidates.sort(function (left, right) {
            return left.id > right.id
        })
    }
});
app.candidateVoted = ko.observable();

app.checkVoteAbility = function() {
    addr = app.web3.eth.accounts[0]; // для примера берём только первый адрес
    console.log("checking "+addr+" address vote ability");
    console.log("checking address to be voter");
    app.ct.voter(addr, function(err, voter) {
        if (err !== null) {
            console.error(err);
            return app.errors.push("Не удалось проверить Вас на возможность избирать")
        }
        app.voter(voter);
        if (!voter) {
            console.log("address is not voter");
            return
        }
        console.log("address is voter");
        console.log("checking address not voted yet");
        app.ct.voted(addr, function(err, voted) {
            if (err !== null) {
                console.error(err);
                return app.errors.push("Не удалось проверить проголосовали ли Вы уже или нет")
            }
            app.voted(voted);
            if (voted) {
                console.log("address is already voted");
                return
            }
            console.log("address is not voted yet");
            app.loadCandidates();
        })
    });
};

app.loadCandidates = function() {
    console.log("loading candidates count");
    app.ct.candidatesCount(function(err, bc) {
        if (err) {
            console.error(err);
            return app.errors.push("Не удалось загрузить количество кандидатов")
        }
        c = bc.toNumber();
        if (c === 0) {
            console.error("unexpected zero candidates count");
            return app.errors.push("Количество кандидатов неожиданно окозалось 0")
        }
        app.candidatesCount(c);
        console.log("candidates count is "+c);
        console.log("loading candidates");
        for (var i = 0; i < c; i++) {
            (function(i) {
                app.ct.candidates(i, function(err, c) {
                    if (err) {
                        console.error("failed to load candidate #"+i+": "+err);
                        return app.errors.push("Не удалось загрузить кандидата №")
                    }
                    console.log("candidate #"+i+" '"+c+"' loaded");
                    app.candidates.push({
                        id: i,
                        name: c,
                        vote: app.vote
                    })
                });
            })(i)
        }
    });
};

app.vote = function(c) {
    console.log("voting for candidate #"+c.id)
    app.ct.vote(c.id, function(err) {
        if (err !== null) {
            console.error(err)
            app.errors.push("Неожиданно не удалось проголосовать")
        }
        console.log("success voted");
        app.candidateVoted(c)
    })
};

ko.applyBindings(app);

app.checkVoteAbility();

    }</script>
</body>
</html>`
