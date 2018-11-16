"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const grpc = require("grpc");
const util = require("util");
const protoLoader = require("@grpc/proto-loader");
const ethers = require("ethers");
let PROTO_PATH = __dirname + '/../proto/moonshard.proto';
class Grpc {
    constructor() {
        this.promisedFuncs = {};
        let packageDefinition = protoLoader.loadSync(PROTO_PATH, {});
        let proto = grpc.loadPackageDefinition(packageDefinition).moonshard;
        this.client = new proto['Moonshard']('localhost:50051', grpc.credentials.createInsecure());
    }
    static GetIdentity() {
        if (!this.identity)
            this.identity = new Grpc();
        return this.identity;
    }
    sign_data(data) {
        if (!this.privKey) {
            throw new Error('private key must be set first');
        }
        let b_m = Buffer.from(data, 'utf-8');
        let sig = this.wal.signingKey.signDigest(ethers.utils.keccak256(b_m));
        let n_sig = ethers.utils.hexlify(ethers.utils.joinSignature(sig));
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
    SetPrivKey(privKey) {
        let wal = new ethers.Wallet(privKey);
        this.wal = wal;
        let addr = wal.address;
        let { privateKey, publicKey } = wal.signingKey.keyPair;
        this.privKey = privateKey;
        this.pubKey = publicKey;
        this.addr = addr;
    }
    CallMethod(method, data) {
        return __awaiter(this, void 0, void 0, function* () {
            if (!this.promisedFuncs[method])
                this.promisedFuncs[method] = util
                    .promisify(this.client[method])
                    .bind(this.client);
            let result = {
                data: {},
                err: null
            };
            try {
                if (data.data) {
                    data.data = JSON.stringify(data.data);
                    data.sign = yield this.sign_data(data.data);
                }
                result.data = yield this.promisedFuncs[method](data);
            }
            catch (e) {
                result.err = e;
            }
            return result;
        });
    }
}
exports.Grpc = Grpc;
//# sourceMappingURL=grpc.js.map