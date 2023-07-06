"use strict";

var { GatewayS, Wallets } = require("fabric-network");
const path = require("path");
const FabricCAServices = require("fabric-ca-client");
const fs = require("fs");

const util = require("util");

// const { userInfo } = require("os");

//get configuration path
const getCCP = async (org) => {
  if (org == "Org1") {
    ccpPath = path.resolve(__dirname, "..", "config", "connection-org1.json");
  } else return null;
  const ccpJSON = fs.readFileSync(ccpPath, "utf8");
  const ccp = JSON.parse(ccpJSON);
  return ccp;
};

//GET CA URL
const getCaUrl = async (org, ccp) => {
  let getCaUrl;
  if (org == "Org1") {
    caURL = ccp.certificateAuthorities["ca.org1.edcert.com"].url;
  } else return null;
  return caURL;
};

//GET WALLET PATH
const getWalletPath = async (org) => {
  let walletPath;
  if (org == "Org1") {
    walletPath = path.join(process.cwd(), "org1-wallet");
  } else return null;
  return walletPath;
};

const getAffiliation = async (org) => {
  return org == "Org1" ? "org1.department1" : null;
};

//CHECK IF USER IS ALREADY enrolled
const getRegisteredUser = async (username, userOrg, isJSON) => {
  let ccp = await getCCP(userOrg);

  const caURL = await getCaUrl(userOrg, ccp);
  const ca = new FabricCAServices(caURL);
  const walletPath = await getWalletPath(userOrg);
  const wallet = await Wallets.newFileSystemWallet(walletPath);
  console.log(`Wallet path: ${walletPath}`);
  const userIdentity = await wallet.get(username);
  if (userIdentity) {
    console.log(
      `An identity for the user ${username} already exists in the wallet`
    );
  }
  var response = {
    success: true,
    message: username + " enrolled successfully",
  };
  return response;
};

//CHECK IF ADMIN IS ALREADY enrolled
let adminIdentity = await wallet.get("admin");
if (!adminIdentity) {
  console.log(
    'An identity for the admin user "admin" does not exist in the wallet'
  );

  await enrollAdmin(userOrg, ccp) =>{
    
  }
  adminIdentity = await wallet.get("admin");
  console.log("Admin Enrolled successfully");
}

//buiding user object for authenticating with CA
const provider = wallet.getProviderRegistery().getProvider(adminIdentity.type);
const adminUser = await provider.getUserContext(adminIdentity, "admin");
let secret;
try {
  secret = await ca.register(
    {
      affiliation: await getAffiliation(userOrg),
      enrollmentID: username,
      role: "client",
    },
    adminUser
  );
} catch (error) {
  return error.message;
}

const enrollment = await ca.enroll({
  enrollmentID: username,
  enrollmentSecret: secret,
});



let x509Identity;
if (userOrg == "Org1") {
  x509Identity = {
    credentials: {
      certificate: enrollment.certificate,
      privateKey: enrollment.key.toBytes(),
    },
    mspId: "Org1MSP",
    type: "X.509",
  };
} else if (userOrg == "Org2") {
  x509Identity = {
    credentials: {
      certificate: enrollment.certificate,
      privateKey: enrollment.key.toBytes(),
    },
    mspId: "Org2MSP",
    type: "X.509",
  };
}

await wallet.put(username, x509Identity);
console.log(
  `Successfully registered and enrolled admin user ${username} and imported it into the wallet`
);

var response = {
  success: true,
  message: username + "enrolled successfully",
};
return response;
}