import * as grpc from 'grpc';
import * as util from 'util';
import * as protoLoader from '@grpc/proto-loader';

import {LocalAddress, CryptoUtils} from 'loom-js';
import * as tweetnacl from 'tweetnacl';

let PROTO_PATH = __dirname + '/../proto/moonshard.proto';

export class Grpc {
    readonly client;
    private static identity:Grpc;
    private promisedFuncs={};
    public addr:string;
    public pubKey:string;
    private privKey:Uint8Array;

    private constructor(){
        let packageDefinition = protoLoader.loadSync(PROTO_PATH,{});
        let proto = grpc.loadPackageDefinition(packageDefinition).moonshard;

        this.client = new proto['Moonshard']('localhost:50051',
            grpc.credentials.createInsecure());
    }

    static GetIdentity(){
        if (!this.identity)
            this.identity=new Grpc();
        return this.identity;
    }

    private sign_data(data) {
        if (!this.privKey){
            throw new Error('private key must be set first');
        }

        let uint_data=Buffer.from(data, 'utf-8');
        let uint8_sig = tweetnacl.sign.detached(
            uint_data, // message as uint8
            this.privKey
        );
        return CryptoUtils.Uint8ArrayToB64(uint8_sig);
    }

    SetPrivKey(privKey:Uint8Array){
        this.privKey=privKey;
        let pub = CryptoUtils.publicKeyFromPrivateKey(this.privKey);
        this.pubKey=CryptoUtils.bytesToHexAddr(pub).toLowerCase();
        this.addr = LocalAddress.fromPublicKey(pub).toString();
    }

    async CallMethod(method:string,data): Promise<{err,data}>{
        if (!this.promisedFuncs[method])
            this.promisedFuncs[method]=util
                .promisify(this.client[method])
                .bind(this.client);
        let result = {
            data:{},
            err:null
        };
        try {
            if (data.data) {
                data.data=JSON.stringify(data.data);
                data.sign=this.sign_data(data.data)
            }
            result.data=await this.promisedFuncs[method](data);
        } catch (e) {
            result.err=e;
        }
        return result;
    }
}