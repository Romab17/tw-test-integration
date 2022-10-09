package main

// #cgo CFLAGS: -I../../include
// #cgo LDFLAGS: -L../../build -L../../build/trezor-crypto -lTrustWalletCore -lprotobuf -lTrezorCrypto -lc++ -lm
// #include <TrustWalletCore/TWHDWallet.h>
// #include <TrustWalletCore/TWPrivateKey.h>
// #include <TrustWalletCore/TWPublicKey.h>
// #include <TrustWalletCore/TWBitcoinScript.h>
// #include <TrustWalletCore/TWAnySigner.h>
// #include <TrustWalletCore/TWTransactionCompiler.h>
import "C"

import (
	"fmt"
	"os"
	"tw/types"
	"unsafe"
)

func main() {
	existingMnemonicStr, mnemonicNotExistErr := os.ReadFile(".mnemonic")

	emptyPassphrase := C.TWStringCreateWithUTF8Bytes(C.CString(""))

	var address unsafe.Pointer

	nullable := types.TWDataCreateWithGoBytes([]byte(""))
	amount := types.TWStringCreateWithGoString("1")
	asset := types.TWStringCreateWithGoString("KSM")
	chainId := types.TWStringCreateWithGoString("40")

	if mnemonicNotExistErr != nil {
		wallet := C.TWHDWalletCreate(128, emptyPassphrase)
		address = C.TWHDWalletGetAddressForCoin(wallet, C.enum_TWCoinType(C.TWCoinTypeKusama))
		mnemonic := C.TWHDWalletMnemonic(wallet)

		os.WriteFile(".mnemonic", []byte(C.GoString(C.TWStringUTF8Bytes(mnemonic))), 0777)

		fmt.Printf("created new wallet.. \n")
	} else {
		mnemonic := C.TWStringCreateWithUTF8Bytes(C.CString(string(existingMnemonicStr)))
		wallet := C.TWHDWalletCreateWithMnemonic(mnemonic, emptyPassphrase)
		address = C.TWHDWalletGetAddressForCoin(wallet, C.enum_TWCoinType(C.TWCoinTypeKusama))

		fmt.Printf("loaded existing wallet.. \n")
	}

	fmt.Printf("wallet address: %s \n", C.GoString(C.TWStringUTF8Bytes(address)))

	valid_address := C.TWCoinTypeValidate(C.enum_TWCoinType(C.TWCoinTypeKusama), address)
	fmt.Printf("Address is valid: ")
	fmt.Println(valid_address)

	input := C.TWTransactionCompilerBuildInput(C.enum_TWCoinType(C.TWCoinTypeKusama), address, address, amount, asset, nullable, chainId)
	C.TWTransactionCompilerPreImageHashes(C.enum_TWCoinType(C.TWCoinTypeKusama), input)
}
