package main

import (
    "encoding/hex"
    "encoding/json"
    "errors"
    "flag"
    "github.com/ethereum/go-ethereum/crypto"
    "gopkg.in/mgo.v2/bson"
    "grpc_service/go_server/dbs"
    "grpc_service/go_server/models"
    "log"
    "net"
    "strings"
    "time"
    
    "context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    pb "grpc_service/proto"
)

//const (
//	port = "50051"
//)

type server struct{}

func getNow() int {
	return int(time.Now().UnixNano() / 1000000)
}

func (s *server) Ping(ctx context.Context, in *pb.PingSetter) (*pb.Success, error) {
	log.Println(in.PubKey)
	pub, _ := DecodeHex(in.PubKey)
	addr := AddrFromPub(pub)
	log.Println("got ping from: " + addr)
	dbs.GetMongo().Instance.C("users").Update(bson.M{"_id": addr}, bson.M{"$set": bson.M{"last_active": getNow()}})
	return &pb.Success{Success: true}, nil
}

func (s *server) GetLatestVersion(ctx context.Context, in *pb.Null) (*pb.LatestVersion, error) {
	var version pb.LatestVersion
	dbs.GetMongo().Instance.C("versions").Find(bson.M{}).Sort("-version").One(&version)
	return &version, nil
}

func (s *server) SetObjData(ctx context.Context, in *pb.ObjSetter) (*pb.Success, error) {
	pub, _ := DecodeHex(in.PubKey)
	valid, _ := CheckSign(pub, in.Data, in.Sign)
	if !valid {
		return nil, errors.New("wrong signature")
	}
	addr := AddrFromPub(pub)
	switch in.Obj {
	case "user":
		byteValue := []byte(in.Data)
		var user models.User
		json.Unmarshal(byteValue, &user)
		if user.Id != addr {
			return nil, errors.New("permission denied")
		}

		// find in dbs
		var exUser models.User
		dbs.GetMongo().Instance.C("users").Find(bson.M{"_id": addr}).One(&exUser)

		// check if exist
		if exUser != (models.User{}) {
			dbs.GetMongo().Instance.C("users").Update(bson.M{"_id": addr}, user)
			return &pb.Success{Success: true}, nil
		}

		err := dbs.GetMongo().Instance.C("users").Insert(user)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &pb.Success{Success: true}, nil
	case "community":
		byteValue := []byte(in.Data)

		// parse data
		var community models.Community
		json.Unmarshal(byteValue, &community)

		// find in dbs
		var exCommunity models.Community
		err := dbs.GetMongo().Instance.C("communities").Find(bson.M{"_id": community.Id}).One(&exCommunity)

		// check if exist
		if exCommunity == (models.Community{}) {
			community.Admin = addr
			dbs.GetMongo().Instance.C("communities").Insert(community)
			return &pb.Success{Success: true}, nil
		}

		// check if setter is admin
		if exCommunity.Admin != addr {
			return nil, errors.New("permission denied")
		}

		// inserting
		err = dbs.GetMongo().Instance.C("communities").Update(bson.M{"_id": community.Id}, community)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &pb.Success{Success: true}, nil
	default:
		return nil, errors.New("wrong obj type")
	}
}

func (s *server) GetObjData(ctx context.Context, in *pb.ObjGetter) (*pb.ObjData, error) {
	var data string
	var upData interface{}
	switch in.Obj {
	case "user":
		var user models.User
		dbs.GetMongo().Instance.C("users").Find(bson.M{"_id": in.Id}).One(&user)
		if user == (models.User{}) {
			return nil, errors.New("no such user")
		}
		upData = user
	case "community":
		var community models.Community
		dbs.GetMongo().Instance.C("communities").Find(bson.M{"_id": in.Id}).One(&community)
		if community == (models.Community{}) {
			return nil, errors.New("no such community")
		}
		upData = community
	}
	dataByte, err := json.Marshal(upData)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	data = string(dataByte)
	return &pb.ObjData{Data: data}, nil
}

func (s *server) GetObjsData(ctx context.Context, in *pb.ObjsGetter) (*pb.ObjData, error) {
	var data string
	var upData interface{}
	switch in.Obj {
	case "user":
		upData = dbs.GetMongo().GetUsers(int(in.Prt), in.Str)
	case "community":
		upData = dbs.GetMongo().GetCommunities(int(in.Prt), in.Str)
	case "all":
		users, communities := dbs.GetMongo().GetAll(int(in.Prt), in.Str)
		upData = struct {
			Users       []models.User
			Communities []models.Community
		}{users, communities}
	}
	dataByte, err := json.Marshal(upData)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	data = string(dataByte)

	return &pb.ObjData{Data: data}, nil
}

func AddrFromPub(pub []byte) string {
	pubE, _ := crypto.UnmarshalPubkey(pub)
	addr := crypto.PubkeyToAddress(*pubE)
	//addr.Hex()
	//hasher := ripemd160.New()
	//hasher.Write(pub)
	//hashBytes := hasher.Sum(nil)
	//address := "0x" + hex.EncodeToString(hashBytes)
	return strings.ToLower(addr.Hex())
}

func DecodeHex(hexStr string) ([]byte, error) {
	decoded, err := hex.DecodeString(hexStr[2:])
	if err != nil {
		log.Println("Decode hexStr err: ", err)
		return nil, err
	} else {
		return decoded, nil
	}
}

func CheckSign(pub []byte, data, sig string) (bool, error) {
	//sig, err := base64.StdEncoding.DecodeString(sigB64)
	//if err != nil {
	//    log.Println("Decode sig err: ",err)
	//    return false, err
	//}
	//dataB:=[]byte(data)
	//sigB,_:=DecodeHex(sig)
	//crypto.VerifySignature(pub,dataB,sigB)
	return true, nil
	//return crypto.VerifySignature(pub,dataB,sigB),nil
}

func main() {
    port := flag.String("p", "50051", "service port")
    host := flag.String("h", "localhost", "service host")
    flag.Parse()
    
	lis, err := net.Listen("tcp", *host+":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMoonshardServer(s, &server{})
	reflection.Register(s)
	log.Println("Server starting at: " + *host+":"+*port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Println("Server down")
}
