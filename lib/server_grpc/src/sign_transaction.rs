#![allow(missing_docs)]

use subxt_signer::sr25519::dev;
use subxt_core::config::{PolkadotConfig};
use subxt_core::config::DefaultExtrinsicParamsBuilder as Params;
use subxt_core::tx;
use subxt_core::utils::H256;
use subxt_core::metadata;
use hex;
use subxt_signer::{ SecretUri, sr25519::Keypair };
use std::str::FromStr;

#[subxt::subxt(runtime_metadata_path = "./metadata.scale")]
pub mod polkadot {}

pub fn voucher_issue(secret_uri_str: &str, bal: &str, upl: bool, dur: u32) -> String{

    // Load metadata
    let metadata_bytes = include_bytes!("../metadata.scale");
    let metadata = metadata::decode_from(&metadata_bytes[..]).unwrap();
    let balance: u128 = bal.parse().expect("Invalid u128");

    let genesis_hash = {
        let h = "d15e9bcfb2e0e00c2c7aa79b6026ded0f7bfe1b813e79266d7683e08e6871625";
        let bytes = hex::decode(h).unwrap();
        H256::from_slice(&bytes)
    };
    // Construct the client state
    let state = tx::ClientState::<PolkadotConfig> {
        metadata,
        genesis_hash,
        runtime_version: tx::RuntimeVersion {
            spec_version: 1810,
            transaction_version: 1,
        },
    };

    let uri = SecretUri::from_str(secret_uri_str).unwrap();
    let keypair = Keypair::from_uri(&uri).unwrap();
    let al =     keypair.public_key().to_account_id();
    println!("{}", al); // returns 5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY for Alice
     
    let call = polkadot::tx().gear_voucher().issue(al, balance, None, upl, dur);
    let params = Params::new().tip(0).nonce(0).build();
    tx::validate(&call, &state.metadata).unwrap();

    let signed_call = tx::create_v4_signed(&call, &state, params)
        .unwrap()
        .sign(&dev::alice());
    println!("Tx: 0x{}", hex::encode(signed_call.encoded()));
    format!("0x{}", hex::encode(signed_call.encoded()))
}
