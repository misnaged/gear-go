package gear_rpc

import "fmt"

type NoArgsMethod int

const (
	MethodRpc NoArgsMethod = iota
	MethodAuthorPendingExtrinsics
	MethodAuthorRotateKeys
	MethodBabeEpochAuthorship
	MethodChainGetFinalizedHead
	MethodGrandpaRoundState
	MethodSystemChain
	MethodSystemChainType
	MethodSystemHealth
	MethodSystemLocalListenAddresses
	MethodSystemLocalPeerId
	MethodSystemName
	MethodSystemNodeRoles
	MethodSystemPeers
	MethodSystemProperties
	MethodSystemReservedPeers
	MethodSystemSyncState
	MethodSystemVersion

	methodUnsupported
)

var NoArgsMethods = [...]string{
	MethodRpc:                        "rpc_methods",
	MethodAuthorPendingExtrinsics:    "author_pendingExtrinsics",
	MethodAuthorRotateKeys:           "author_rotateKeys",
	MethodBabeEpochAuthorship:        "babe_epochAuthorship",
	MethodChainGetFinalizedHead:      "chain_getFinalizedHead",
	MethodGrandpaRoundState:          "grandpa_roundState",
	MethodSystemChain:                "system_chain",
	MethodSystemChainType:            "system_chainType",
	MethodSystemHealth:               "system_health",
	MethodSystemLocalListenAddresses: "system_localListenAddresses",
	MethodSystemLocalPeerId:          "system_localPeerId",
	MethodSystemName:                 "system_name",
	MethodSystemNodeRoles:            "system_nodeRoles",
	MethodSystemPeers:                "system_peers",
	MethodSystemProperties:           "system_properties",
	MethodSystemReservedPeers:        "system_reservedPeers",
	MethodSystemSyncState:            "system_syncState",
	MethodSystemVersion:              "system_version",
}

func (s NoArgsMethod) String() string {
	return NoArgsMethods[s]
}
func NoArgMethodFromString(s string) NoArgsMethod {
	for i, r := range NoArgsMethods {
		if s == r {
			return NoArgsMethod(i)
		}
	}
	return methodUnsupported
}

// NoArgMethodFromStringE return new NoArgsMethod enum
// from given string or return an error
func NoArgMethodFromStringE(s string) (NoArgsMethod, error) {
	for i, r := range NoArgsMethods {
		if s == r {
			return NoArgsMethod(i), nil
		}
	}
	return methodUnsupported, fmt.Errorf("invalid method value %q", s)
}
