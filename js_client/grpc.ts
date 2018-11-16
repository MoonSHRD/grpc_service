import * as grpc from 'grpc';
import * as util from 'util';
import * as protoLoader from '@grpc/proto-loader';
import ethers = require('ethers');

import {LocalAddress, CryptoUtils} from 'loom-js';
import * as tweetnacl from 'tweetnacl';

let PROTO_PATH = __dirname + '/../proto/moonshard.proto';

export class Grpc {
    readonly client;
    private static identity:Grpc;
    private promisedFuncs={};
    public addr:string;
    public pubKey:string;
    private privKey:string;
    private wal:any;

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

        let b_m=Buffer.from(data,'utf-8');
        let sig=this.wal.signingKey.signDigest(ethers.utils.keccak256(b_m));
        let n_sig=ethers.utils.hexlify(ethers.utils.joinSignature(sig));
        return n_sig;
        // let uint_data=Buffer.from(data, 'utf-8');
        // let uint8_sig = tweetnacl.sign.detached(
        //     Buffer.from(data, 'utf-8'), // message as uint8
        //     Buffer.from(this.privKey.substr(2), "hex")
        // );
        // Buffer.from(uint8_sig.buffer, uint8_sig.byteOffset, uint8_sig.byteLength)
        //     .toString('hex')
        //     .toUpperCase();
        // return uint8_sig.toString('hex');
    }

    SetPrivKey(privKey:string){
        let wal=new ethers.Wallet(privKey);
        this.wal=wal;
        let addr = wal.address;
        let {privateKey,publicKey}=wal.signingKey.keyPair;
        this.privKey=privateKey;
        this.pubKey=publicKey;
        this.addr = addr;
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
                data.sign=await this.sign_data(data.data);
            }
            result.data=await this.promisedFuncs[method](data);
        } catch (e) {
            result.err=e;
        }
        return result;
    }
}