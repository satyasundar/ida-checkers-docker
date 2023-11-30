package app

import (
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/ignite-hq/cli/ignite/pkg/cosmoscmd"
	appparams "github.com/satya/checkers/app/params"
)

func MakeTestEncodingConfig() cosmoscmd.EncodingConfig {
	encodingConfig := appparams.MakeTestEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
