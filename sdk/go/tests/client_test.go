package tests

import (
    "encoding/hex"
    "testing"
    . "sawtooth_sdk/client"
)

var (
  data = []byte{0x01, 0x02, 0x03}
  PEMSTR = `-----BEGIN EC PRIVATE KEY-----
MHQCAQEEIISGREvlByLRnbhovpK8wSd5hnymtY8hdQCOvZQ473CpoAcGBSuBBAAK
oUQDQgAEWC6TyM1jpYu3f/GGIuktDk4nM1qyOf9PEPHkRkN8zK2HxxNwDi+yN3hR
8Ag+VeTwbRRZOlBdFBsgPxz3/864hw==
-----END EC PRIVATE KEY-----
`
  PEMSTRPRIV = "8486444be50722d19db868be92bcc12779867ca6b58f2175008ebd9438ef70a9"
  ENCPEM = `-----BEGIN EC PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: AES-128-CBC,23CDF282F2217A9334A2413D78DAE04C

PQy89wdLsayP/FG68wgmL1EdlI3S5pN8ibCFrnp5OAtVNrYUD/TH9DMYVmRCNUB4
e+vXoQzd1IysjoFpV21zajSAxCmcbU4CGCDEea3GPwirOSE0ZjPHPp15IkRuGFYm
L/8e9mXvEQPAmBC0NMiltnk4/26iN7hB1QxSQQwy/Zc=
-----END EC PRIVATE KEY-----
`
  ENCPEMPRIV = "2cc32bc33935a5dbad8118abc63dfb627bb91a98d5e6310f5d60f5d65f6adb2f"
  PEMPUBSTR = "03582e93c8cd63a58bb77ff18622e92d0e4e27335ab239ff4f10f1e446437cccad"
  ENCPEMPUB = "0257510b4718fd79b21dee3173ffb48ab9a668a35a377be7b7dc432243a940c510"
  WIFSTR = "5J7bEeWs14sKkz7yVHfVc2FXKfBe6Hb5oNZxxTKqKZCgjbDTuUj"
  PUBSTR = "035e1de3048a62f9f478440a22fd7655b80f0aac997be963b119ac54b3bfdea3b7"
  SIGSTR = "0062bc154dca72472e66062c4539c8befb2680d79d59b3cc539dd182ff36072b199adc1118db5fc1884d50cdec9d31a2356af03175439ccb841c7b0e3ae83297"
  ENCDED = "0acc0a0aca020a423033356531646533303438613632663966343738343430613232666437363535623830663061616339393762653936336231313961633534623362666465613362371280013266363339316161336331643930633436653863353131623766363832303965386232363131376363303639333564623839383538636633663863643232373032646339363836663434346232613263356632346463396530653336353363393566313666373264396438343863383731356562663561623630323264343839128001306634626631653631636337613830666566373930313733356662353264633562323333323333333463373166366631656331353838343061643865306563353739316663333332626538303435623830656130316539323461663434336135306334363231653036323535366432616565613332333235636338303434643712800133653563326432306365653335643431643366333062666434346661376463363130633337346230313033313666633138646365623261316538326531636161336137316362346334366530393363336464353436306235396133393636653533383861393832623162353163306164393833363461633966613837323932301abb030ab0020a423033356531646533303438613632663966343738343430613232666437363535623830663061616339393762653936336231313961633534623362666465613362371a0361626322033132332a0364656632033132333a03646566420a6d79656e636f64696e674a80013237383634636335323139613935316137613665353262386338646464663639383164303938646131363538643936323538633837306232633838646662636235313834316165613137326132386261666136613739373331313635353834363737303636303435633935396564306639393239363838643034646566633239524230333565316465333034386136326639663437383434306132326664373635356238306630616163393937626539363362313139616335346233626664656133623712800132663633393161613363316439306334366538633531316237663638323039653862323631313763633036393335646238393835386366336638636432323730326463393638366634343462326132633566323464633965306533363533633935663136663732643964383438633837313565626635616236303232643438391a030102031abb030ab0020a423033356531646533303438613632663966343738343430613232666437363535623830663061616339393762653936336231313961633534623362666465613362371a0361626322033132332a0364656632033435363a03676869420a6d79656e636f64696e674a80013237383634636335323139613935316137613665353262386338646464663639383164303938646131363538643936323538633837306232633838646662636235313834316165613137326132386261666136613739373331313635353834363737303636303435633935396564306639393239363838643034646566633239524230333565316465333034386136326639663437383434306132326664373635356238306630616163393937626539363362313139616335346233626664656133623712800130663462663165363163633761383066656637393031373335666235326463356232333332333333346337316636663165633135383834306164386530656335373931666333333262653830343562383065613031653932346166343433613530633436323165303632353536643261656561333233323563633830343464371a03010203"
)

func TestSigning(t *testing.T) {
    priv := GenPrivKey()
    pub := GenPubKey(priv)
    sig := Sign(data, priv)
    if !Verify(data, sig, pub) {
        t.Error(
            "Couldn't verify generated signature",
            priv, pub, sig,
        )
    }
}

func TestEncoding(t *testing.T) {
    priv, err := WifToPriv(WIFSTR)
    if err != nil {
        t.Error("Failed to load WIF key")
    }
    if PrivToWif(priv) != WIFSTR {
        t.Error("Private key is different after encoding/decoding")
    }
    if hex.EncodeToString(GenPubKey(priv)) != PUBSTR {
        t.Error("Public key doesn't match expected. Got", GenPubKey(priv))
    }
    sigstr := hex.EncodeToString(Sign(data, priv))
    if sigstr != SIGSTR {
        t.Error("Signature doesn't match expected. Got", sigstr)
    }
}

func TestPemLoader(t *testing.T) {
    // Load the keys
    priv, err := PemToPriv(PEMSTR, "")
    if err != nil {
        t.Error("Failed to load unencrypted PEM key")
    }
    epriv, err := PemToPriv(ENCPEM, "password")
    if err != nil {
        t.Error("Failed to load encrypted PEM key")
    }
    // Test that they match expected
    if hex.EncodeToString(priv) != PEMSTRPRIV {
        t.Error("Failed to parse unencrypted PEM key")
    }
    if hex.EncodeToString(epriv) != ENCPEMPRIV {
        t.Error("Failed to parse encrypted PEM key")
    }
    // Test that the correct public keys are generated
    pub := hex.EncodeToString(GenPubKey(priv))
    epub := hex.EncodeToString(GenPubKey(epriv))
    if pub != PEMPUBSTR {
      t.Error("Failed to generate correct public key from unencrypted PEM key")
    }
    if epub != ENCPEMPUB {
      t.Error("Failed to generate correct public key from encrypted PEM key")
    }
}

func TestEncoder(t *testing.T) {
    priv, _ := WifToPriv(WIFSTR)

    encoder := NewEncoder(priv, TransactionParams{
        FamilyName: "abc",
        FamilyVersion: "123",
        PayloadEncoding: "myencoding",
        Inputs: []string{"def"},
    })

    txn1 := encoder.NewTransaction(data, TransactionParams{
        Nonce: "123",
        Outputs: []string{"def"},
    })

    pubstr := hex.EncodeToString(GenPubKey(priv))
    txn2 := encoder.NewTransaction(data, TransactionParams{
        Nonce: "456",
        Outputs: []string{"ghi"},
        BatcherPubkey: pubstr,
    })

    // Test serialization
    txns, err := ParseTransactions(SerializeTransactions([]*Transaction{txn1, txn2}))
    if err != nil {
        t.Error(err)
    }

    batch := encoder.NewBatch(txns)

    // Test serialization
    batches, err := ParseBatches(SerializeBatches([]*Batch{batch}))
    if err != nil {
        t.Error(err)
    }
    data := SerializeBatches(batches)
    datastr := hex.EncodeToString(data)

    expected := ENCDED

    if datastr != expected {
        t.Error("Did not correctly encode batch. Got", datastr)
    }
}
