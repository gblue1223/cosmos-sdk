# 튜토리얼: Cosmos-SDK 애플리케이션을 코딩하는 방법

이 튜토리얼에서는 Cosmos-SDK를 사용하여 애플리케이션을 코딩하는 기본 사항을 학습합니다. 텐더민트 (Tendermint)와 코스모스 생태계 (Cosmos Ecosystem)에 대한 소개부터 시작하여 텐더 민트 (Tendermint) 소프트웨어와 Cosmos-SDK 프레임워크에 대한 개괄적인 개요를 살펴 보겠습니다. 그 다음에는 코드에 대해 알아보고 Cosmos-SDK 제일 윗단에서 간단한 애플리케이션을 작성하는 방법을 안내합니다.

## 텐더민트와 코스모스

블록 체인은 세 가지 개념적 레이어로 나눌 수 있습니다:

- **네트워킹:** 트랜잭션 전파에 책임이 있습니다.
- **합의:** 다음 트랜잭션 집합의 처리에 동의하기 위해 검증인 노드들을 활성화 합니다 (예: 블록체인에 트랜잭션 블록 추가).
- **애플리케이션:** 트랜잭션 집합에 주어진 상태를 업데이트하는 작업을 담당합니다. 즉, 트랜잭션을 처리합니다.

*네트워킹* 계층은 각 노드가 트랜잭션을 수신 하는지 확인합니다. *합의* 계층은 각 노드가 자신의 로컬 상태를 수정하기 위해 동일한 트랜잭션에 동의하는지 확인합니다. *애플리케이션* 레이어는 트랜잭션을 처리합니다. 트랜잭션과 상태가 주어지면 애플리케이션은 새로운 상태를 반환합니다. Bitcoin을 예를 들면, 애플리케이션 상태는 각 계정에 대한 원장 또는 잔액 목록입니다 (실제로는 UTXO의 목록으로 Unspent Transaction Output의 약자이지만 단순화를 위해 잔액이라고합시다). 그리고 트랜잭션은 이러한 잔액 목록을 변경하여 애플리케이션의 상태를 수정합니다. Ethereum의 경우 애플리케이션은 가상 머신입니다. 각 트랜잭션은 이 가상 머신을 통해 진행되며 애플리케이션의 상태가 트랜잭션 내부에서 호출된 스마트 계약에 따라 수정됩니다. 

텐더민트 이전에는, 블록 체인을 만들려면 세 개의 레이어를 처음부터 모두 만들어야했습니다. 대부분의 개발자가 Bitcoin 코드베이스를 포크하거나 복제하는 것을 선호했지만 Bitcoin의 프로토콜의 한계에 의해 제약을받기 때문에 실증나는 작업이었습니다. Ethereum Virtual Machine (EVM)은 이 문제를 해결하고 스마트 계약을 통해 사용자 정의 가능한 논리를 실행할 수 있게 하여 분산된 애플리케이션 개발을 단순화하도록 설계되었습니다. 그러나 블록 체인 자체의 한계 (상호 운용성, 확장 성 및 사용자 정의)를 해결하지 못했습니다. Go-Ethereum은 Bitcoin의 코드베이스처럼 하드 포크하기가 어려운 매우 획일적인 기술 스택입니다. 

Tendermint는 이러한 문제를 해결하고 개발자에게 대안을 제공하기 위해 고안되었습니다. Tendermint의 목표는 빌드 하려는 모든 애플리케이션 개발자를 지원하는 일반 엔진으로서 블록 체인의 *네트워킹* 및 *합의* 레이어를 제공하는 것입니다. Tendermint를 사용하면 개발자는 *애플리케이션* 레이어에만 집중해야하므로 수백시간의 작업시간과 값 비싼 개발 환경을 절약 할 수 있습니다. 참고로 Tendermint는 Tendermint Core 엔진 내에서 사용 되는 비잔틴 장애 허용 합의 알고리즘을 나타냅니다.

Tendermint는, [ABCI] (https://github.com/tendermint/abci)(Application-BlockChain Interface) 라 불리는 소켓 프로토콜을 통해 블록체인 엔진인 Tendermint Core (*네트워킹* 및 *합의* 레이어)를 *애플리케이션* 레이어에 연결합니다. 개발자는 Tendermint Core 엔진에서 실행되는 ABCI 지원 애플리케이션을 작성하기 위해 몇 가지 메시지만 구현 하면 됩니다. ABCI는 언어에 구속력이 없습니다. 즉, 개발자는 어떤 프로그래밍 언어를 이용해서 블록 체인의 *애플리케이션* 파트를 구현할 수 있습니다. Tendermint Core 엔진 위에 구축하면 다음과 같은 이점도 있습니다.

- ** 공용 또는 사설 블록 체인이 가능합니다.** 개발자는 Tendermint Core 위에 권한이 부여 된 (사설) 권한이 없는 (공용) 모든 블록 체인 애플리케이션을 배포 할 수 있습니다. 
-** 성능.** 텐더 민트 코어는 짧은 시간 간격으로 많은 수의 트랜잭션을 처리 할 수있는 최첨단 블록 체인 컨센서스 엔진입니다. Tendermint Core의 블록 시간은 1 초 정도로 낮을 수 있으며 해당 기간에 수천 회의 트랜잭션을 처리 할 수 ​​있습니다.
- ** 즉각적인 완결성.** Tendermint 컨센서스 알고리즘의 속성은 즉각적인 완결성입니다. 즉 유효성 검사기의 1/3 미만이 악성 코드 (비잔틴) 인 한 포크는 생성되지 않습니다. 사용자는 블록이 생성되는 즉시 트랜잭션이 완결되었는지 확인할 수 있습니다.
-**보안.** Tendermint Core의 합의는 단지 장애를 허용하는 것이 아니라, 책임감있는 최적의 비잔틴 장애 허용 (BFT) 에 있습니다. 만일 블록 체인이 포크되면, 책임을 밝힐 수있는 방법이 있습니다.
- **라이트 클라이언트 지원**. 텐더민트는 내장된 라이트 클라이언트를 제공합니다.

하지만 가장 중요한 점은 텐더민트는 [Inter-Blockchain Communication Protocol] (https://github.com/cosmos/cosmos-sdk/tree/develop/docs/spec/ibc) (IBC)과 호환된다는 것입니다. 즉, 공용이든 사설이든 텐더민트 기반 블록 체인은 본질적으로 코스모스 생태계와 연결되어 생태계의 다른 블록 체인과 토큰을 안전하게 교환 할 수 있습니다. IBC와 코스모스를 통한 상호 운용성의 혜택은 텐더민트 체인의 자주권을 보호합니다. 비 텐더민트 체인은 IBC 어댑터 또는 페그존을 통해 코스모스에 연결할 수도 있지만 이 내용은 이 문서의 범위를 벗어납니다.

코스모스 생태계에 대한 자세한 내용은 [이 글] (https://blog.cosmos.network/understanding-the-value-proposition-of-cosmos-ecaef63350d)을 참조하십시오.


## 코스모스SDK 소개

텐더민트 기반 블록 체인을 개발한다는 것은 애플리케이션 (즉, 상태 머신) 만 코딩 하면 된다는 것을 의미합니다. 하지만 그 자체로는 다소 어려울 수 있습니다. 이것이 바로 Cosmos-SDK가 존재하는 이유입니다.

[Cosmos-SDK] (https://github.com/cosmos/cosmos-sdk)는 Cosmos 허브와 같은 다중 자산 Proof-of-Stake (PoS) 블록 체인이자 Proof-Of -Authority (PoA) 블록 체인을 구축하기위한 플랫폼.

Cosmos-SDK의 목표는 개발자가 일반적인 블록 체인 기능을 다시 만들 필요 없이 Cosmos 네트워크 내에서 상호 운용 가능한 맞춤 블록 첸인 애플리케이션을 쉽게 만들 수있게하고, Tendermint ABCI 애플리케이션을 작성하는 복잡성을 제거하는 것입니다. 우리는 Tendermint 위에 안전한 블록 체인 어플리케이션을 구축하기 위해 npm과 유사한 프레임워크의 SDK를 구상하고 있습니다.

SDK는 설계 측면에서 유연성과 보안성을 최대한 신경쓰고 있습니다. 프레임워크는 애플리케이션이 원하는대로 요소를 혼합하고 일치시킬 수있는 모듈 실행 스택 위주로 설계되었습니다. 또한 모든 모듈은 보다 강력한 애플리케이션 보안을 위해 샌드박스화되어 있습니다.

이것은 두 가지 주요 원칙에 기반합니다:

- ** 합성성:** 누구나 Cosmos-SDK 용 모듈을 만들 수 있으며 이미 구축 된 모듈을 통합하는 것은 블록체인 어플리케이션으로 가져 오는 것 만큼 간단합니다.

- **기능들:** SDK는 기능 기반 보안에서 영감을 얻었으며, 수년간 블록 체인 상태 머신과의 씨름을 통해 영향을 끼쳤습니다. 대부분의 개발자는 자체 모듈을 만들 때 다른 타사 모듈에 액세스 해야합니다. Cosmos-SDK는 개방형 프레임워크이며 이러한 모듈 중 일부는 악의적이라고 가정하기 때문에, 객체 기능 (OCAPS) 기반 원칙을 사용하여 SDK를 설계했습니다. 실제로, 이는 각 모듈이 다른 모듈에 대한 액세스 제어 목록을 유지하는 대신, 각 모듈은 사전 정의 된 기능 세트를 부여하기 위해 다른 모듈에 전달할 수 있는 키퍼라는 특수 객체를 구현한다는 것을 의미합니다. 예를 들어 모듈 A의 키퍼 인스턴스가 모듈 B로 전달되면, 후자는 모듈 A의 제한된 함수 집합을 호출 할 수 있습니다. 각 키퍼의 기능은 모듈 개발자가 정의하며 각 타사 모듈에 전달되는 기능을 기반으로 타사 모듈의 외부 코드 안전성을 이해하고 감사하는 것은 개발자의 임무입니다. '기능'에 대한 더 자세한 내용은 [이 글](http://habitatchronicles.com/2017/05/what-are-capabilities/)을 참조하십시오.

*참고: 현재 Cosmos-SDK는 Golang에만 존재합니다. 즉, 개발자는 Golang으로만 SDK 모듈을 개발할 수 있습니다. 앞으로, SDK는 다른 프로그래밍 언어로 구현 될 것이라 기대됩니다. Tendermint 팀에서 자금을 지원하는 것이 결국에는 가능할 것입니다.*

## Tendermint 및 ABCI에 대한 알림

Cosmos-SDK는 블록 체인의 *애플리케이션* 레이어를 개발하기위한 프레임워크입니다. 이 애플리케이션은 [Application-Blockchain Interface] (https://github.com/tendermint/abci)의 약자 인 ABCI라는 간단한 프로토콜을 지원하는 모든 합의 엔진 (*합의* + *네트워킹* 레이어)에 연결될 수 있습니다.

Tendermint Core는 Cosmos-SDK 위에 구축된 기본 컨센서스 엔진입니다. *애플리케이션* 및 *합의 엔진* 각자의 책임을 잘 이해하는 것이 중요합니다.

*합의 엔진*의 책임:
- 트랜잭션 전파
- 유효한 트랜잭션의 순서에 대한 동의.

*애플리케이션*의 책임:
- 트랜잭션 생성
- 트랜잭션이 유효한지 확인.
- 트랜잭션 처리 (상태 변경 포함)

*합의 엔진*은 각 블록에 대해 제공된 유효성 검증인 집합에 대한 지식을 갖고 있지만 유효성 검사기 집합 변경을 발동하는 것은 *애플리케이션*의 책임이라는 점을 강조 할 가치가 있습니다. 이것이 Cosmos-SDK와 Tendermint를 통해 **공개 체인과 사설 체인**을 모두 구축 할 수 있는 이유입니다. 체인은 유효성 검증인의 설정 변경을 제어하는 ​​애플리케이션 수준에서 정의 된 규칙에 따라 공개 또는 사설이됩니다.

ABCI는 *합의 엔진*과 *애플리케이션* 사이의 연결을 설정합니다. 본질적으로, 핵심은 두 가지 메시지입니다:

-`CheckTx`: 트랜잭션이 유효한지 어플리케이션에게 물어 봅니다. 검증인의 노드가 트랜잭션을 받으면 검증인의 노드는 `CheckTx`를 실행합니다. 트랜잭션이 유효하면 메모리풀에 추가됩니다.
-`DeliverTx`: 애플리케이션에게 트랜잭션을 처리하고 상태를 업데이트하도록 요청합니다.

*합의 엔진*과 *애플리케이션*이 서로 어떻게 상호 작용 하는지에 대한 개괄적인 개요를 설명하겠습니다.

- 항상, 검증인 노드의 합의 엔진 (Tendermint Core)이 트랜잭션을 수신 할 때마다 CheckTx를 통해 애플리케이션에 전달하여 유효성을 검사합니다. 유효하면 트랜잭션이 메모리풀에 추가됩니다.
- 우리가 블록 N에 있다고합시다. 유효성 검증인 세트 V가 있습니다. 다음 블록의 제안자는 *합의 엔진*에 의해 V에서 선택됩니다. 제안자는 메모리풀에서 유효한 트랜잭션을 수집하여 새 블록을 만듭니다. 그런 다음 블록은 다른 검증인들에게 알려지고 서명/커밋 됩니다. V의 2/3 이상이 *사전커밋*에 서명하면 블록은 블록 N+1이됩니다 (합의 알고리즘에 대한 자세한 설명은 [여기] (https://github.com/tendermint/tendermint/wiki/Byzantine-Consensus-Algorithm)를 클릭하십시오.
- 블록 N+1이 V의 2/3 이상에 의해 서명 될 때, 그것은 full-node로 알려지게 됩니다. full-node가 해당 블록을 수신하면 유효성을 확인합니다. 블록이 V의 2/3 이상 유효한 서명들을 보유하고 블록의 모든 트랜잭션들이 유효한 경우 블록은 유효합니다. *합의 엔진*은 트랜잭션의 유효성을 검사하기 위해 'DeliverTx'를 통해 애플리케이션으로 전송합니다. 각 트랜잭션 이후에 'DeliverTx`는 트랜잭션이 유효하면 새로운 상태를 반환합니다. 블록이 끝나면 최종 상태가 확약됩니다. 물론, 이것은 블록 내의 트랜잭션 순서가 중요하다는 것을 의미합니다.

## SDK-app의 아키텍처

Cosmos-SDK는 Tendermint 기반 블록 체인 애플리케이션을 위한 기본 템플릿을 제공합니다. 이 템플릿은 [여기] (https://github.com/cosmos/cosmos-sdk)에서 찾을 수 있습니다.

본질적으로, 블록 체인 애플리케이션은 단순히 복제된 상태 머신입니다. 상태 (예: 암호화폐, 각 계정이 보유하는 동전 수) 및 상태 전이를 트리거하는 트랜잭션이 있습니다. 애플리케이션 개발자는 상태, 트랜잭션 유형 및 다른 트랜잭션이 상태를 수정하는 방법을 정의합니다.

### 모듈성

Cosmos-SDK는 모듈 기반 프레임워크입니다. 각 모듈은 그 자체로 작은 상태 기계이며 다른 모듈과 쉽게 결합되어 일관된 애플리케이션을 생성 할 수 있습니다. 즉, 모듈은 전역 상태와 트랜잭션 유형의 하위 섹션을 정의합니다. 그런 다음, 각 유형에 따라 올바른 모듈로 트랜잭션을 라우팅하는 것은 루트 애플리케이션의 일입니다. 이 과정을 이해하기 위해 상태 머신의 단순화된 일반적인 주기를 살펴 보겠습니다.

Tendermint Core 엔진에서 트랜잭션을 수신하면, 다음은 *애플리케이션*이 수행하는 작업입니다:

1. 트랜잭션을 디코드하고 메시지를 받음
2.`Msg.Type()` 메소드를 사용하여 메시지를 적절한 모듈로 라우팅
3. 모듈에서 트랜잭션을 실행. 트랜잭션이 유효하면 상태를 수정.
4. 새 상태 또는 오류 메시지를 반환

단계 1, 2 및 4는 루트 애플리케이션에서 처리합니다. 3 단계는 해당 모듈에서 처리합니다.

### SDK 구성 요소

이를 염두에두고 SDK의 중요한 디렉토리를 살펴 보겠습니다:

-`baseapp`: 기본 애플리케이션을위한 템플릿을 정의. 기본적으로 Cosmos-SDK 응용 프록로그램이 기본 Tendermint 노드와 통신 할 수 있도록 ABCI 프로토콜을 구현합니다.
-`client`: 애플리케이션과 상호 작용하는 커맨드라인 인터페이스
-`server`: 노드와 통신하기 위한 REST 서버
-`examples`: `baseapp`과 모듈을 기반으로 작동하는 애플리케이션을 만드는 법에 대한 예제가 들어 있음
-`store`: 멀티 스토어에 대한 코드를 포함. 멀티 스토어는 기본적으로 사용자의 상태입니다. 각 모듈은 멀티 스토어에서 원하는 수의 KVStore들을 생성 할 수 있습니다. 해당 `키퍼`를 사용하여 각 상점에 대한 액세스 권한을 올바르게 처리하도록 주의하십시오.
-`types`: SDK 기반의 모든 애플리케이션에 필요한 공통 유형.
-`x`: 모듈이 있는 곳. 이 디렉토리에는 이미 빌드된 모듈이 모두 있습니다. 이러한 모듈 중 하나를 사용하려면, 애플리케이션에서 적절하게 가져와야합니다. [App - 모든 것을 하나로 연결] 섹션에서 우리는 사용방법을 보게 될 것입니다.

### 입문자들을 위한 코드런

#### KV스토어

KV스토어는 SDK 애플리케이션의 기본 지속성 계층을 제공합니다.

```go
type KVStore interface {
    Store

    // Get 키가 존재하지 않으면 nil을 반환합니다. nil 키면 패닉실행.
    Get(key []byte) []byte

    // Has 키가 있는지 검사합니다. nil 키면 패닉실행.
    Has(key []byte) bool

    // Set 키를 설정합니다. nil 키면 패닉실행.
    Set(key, value []byte)

    // Delete 키를 삭제합니다. nil 키면 패닉실행.
    Delete(key []byte)

    // 오름차순으로 키 도메인에 대한 반복자. End는 배타적임.
    // Start는 end보다 작아야합니다. 그렇지 않으면 iterator가 유효하지 않습니다.
    // CONTRACT: 반복자가 있는 동안은 도메인내에서 쓰기가 발생하지 않습니다.
    Iterator(start, end []byte) Iterator

    // 오름차순으로 키 도메인에 대한 반복자. End는 배타적임.
    // Start는 end보다 작아야합니다. 그렇지 않으면 iterator가 유효하지 않습니다.
    // CONTRACT: 반복자가 있는 동안은 도메인내에서 쓰기가 발생하지 않습니다.
    ReverseIterator(start, end []byte) Iterator

    // TODO 아직 구현되지 않았습니다.
    // CreateSubKVStore(key *storeKey) (KVStore, error)

    // TODO 아직 구현되지 않았습니다.
    // GetSubKVStore(key *storeKey) KVStore
 }
```

계정에 하나, IBC에 하나 등과 같이, 여러 개의 KVStore를 애플리케이션에 마운트 할 수 있습니다.

```go
 app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keyStake, app.keySlashing)
```

요청이있는 경우, KVStore의 구현은 각 쿼리에 대해 Merkle 증명을 제공 할 책임이 있습니다.

```go
 func (st *iavlStore) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
```

KVStore는 지속성 레벨에서 트랜잭션을 제공 할 수 있도록 캐시래핑 될 수 있습니니다 (iterator에서도 잘 지원됩니다). 이 기능은 "AnteHandler"가 트랜잭션에 대한 모든 관련 비용을 공제 한 후 트랜잭션 처리를 위한 트랜잭션 격리 계층을 제공하는 데 사용됩니다. 캐시래핑은 블록 체인을위한 가상 머신 또는 스크립팅 환경을 구현할 때 유용 할 수 있습니다.

#### go-amino

Cosmos-SDK는 Go 형식을 Protobuf3 호환 바이트로 직렬화 및 비직렬화 하기 위해 [go-amino] (https://github.com/cosmos/cosmos-sdk/blob/96451b55fff107511a65bf930b81fb12bed133a1/examples/basecoin/app/app.go#L97-L111)를 광범위하게 사용합니다.

Go-amino (예 : https://github.com/golang/protobuf)는 리플렉션을 사용하여 Go 객체를 인코딩/디코딩합니다.  이를 통해 SDK 개발자는 Proto3에 대한 별도의 스키마를 유지할 필요없이 Go에서 데이터 구조를 정의하는데 집중할 수 있습니다. 또한, Amino는 인터페이스 및 실제 타입을 위해 네이티브 지원을 통해 Proto3을 확장합니다.

예를 들어, Cosmos-SDK의 `x/auth` 패키지는 PubKey 인터페이스를`tendermint/go-crypto` 에서 가져옵니다. PubKey 구현은 _Ed25519_ 및 _Secp256k1_에 대한 구현을 포함합니다.  각 `auth.BaseAccount`에는 PubKey가 있습니다.

```go
 // BaseAccount - 기본 계정 구조.
 // AppAccount에 이것을 포함시켜 확장합니다.
 // 예제를 보려면 examples/basecoin/types/account.go를 참조하십시오.
 type BaseAccount struct {
    Address  sdk.Address   `json:"address"`
    Coins    sdk.Coins     `json:"coins"`
    PubKey   crypto.PubKey `json:"public_key"`
    Sequence int64         `json:"sequence"`
 }
```

Amino는 인터페이스에 등록 된 실제 값을 기반으로 각 인터페이스 값에 대해 디코딩할 실제 타입을 알고 있습니다.

예를 들어, `Basecoin` 예제 애플리케이션은 _Ed25519_와 _Secp256k1_ 키에 대해 알고 있는데, 그 이유는 아래의 앱의 'codec'에 의해 등록 되었기 때문입니다.

```go
wire.RegisterCrypto(cdc) // 암호를 등록.
```

Go-Amino에 대한 자세한 내용은 https://github.com/tendermint/go-amino를 참조하십시오.

#### Keys, Keepers, 그리고 Mappers

Cosmos-SDK는 전체 애플리케이션을 구성하기 위해 함께 포함될 수 있는 라이브러리의 생태계를 활성화할 수 있도록 설계되었습니다. 이 생태계를보다 안전하게 만들기 위해 최소-권한 원칙에 따라 디자인 패턴을 개발했습니다.

`Mappers`와 `Keepers`는 컨텍스트를 통해 KVStore 들에 대한 액세스를 제공합니다. 이 둘의 유일한 차이점은 `Mapper`가 낮은 수준의 API를 제공하기 때문에 일반적으로 `Keepr`는 다른 `Keepers` 및 `Mappers`에 대한 참조를 보유 할 수 있지만 그 반대는 아닙니다.

`Mappers`와 `Keepers`는 직접 KVStore들을 참조하지 않습니다.  그들은 오직 _key_ (아래의 `sdk.StoreKey`) 만을 가지고 있습니다:

```go
type AccountMapper struct {

    // 컨텍스트에서 KVStore에 액세스하는 데 사용된 (노출되지 않은) 키입니다.
    key sdk.StoreKey

    // 계정 실제 타입.
    proto Account

    // 계정의 바이너리 인코딩/디코딩을 위한 와이어 코덱.
    cdc *wire.Codec
 }
```

이런 식으로, 당신은 `app.go` 파일에 모든 것을 연결하고 어떤 컴포넌트가 어떤 KVStore 들과 다른 컴포넌트에 접근 할 수 있는지를 볼 수 있습니다.

```go
// accountMapper를 정의합니다.
 app.accountMapper = auth.NewAccountMapper(
    cdc,
    app.keyAccount, // 대상 KVStore
    &types.AppAccount{}, // prototype
 )
```

나중에 트랜잭션을 실행하는 동안 (예: 블록 커밋 후에 ABCI `DeliverTx`를 통해) 컨텍스트가 첫 번째 인수로 전달됩니다.  컨텍스트에는 모든 관련 KVStore 들에 대한 참조가 포함되지만 관련 키를 보유한 경우에만 액세스 할 수 있습니다.

```go
 // sdk.AccountMapper를 구현합니다.
 func (am AccountMapper) GetAccount(ctx sdk.Context, addr sdk.Address) Account {
    store := ctx.KVStore(am.key)
    bz := store.Get(addr)
    if bz == nil {
        return nil
    }
    acc := am.decodeAccount(bz)
    return acc
 }
```

`Mappers`와 `Keepers`는 앱 초기화시 KVStore가 알려지지 않아 KVStore에 대한 직접적인 참조를 가질 수 없습니다.  KVStore는 모든 커밋 된 트랜잭션 (ABCI `DeliverTx`를 통해) 및 메모리풀 확인 트랜잭션 (ABCI `CheckTx`를 통해)에 대해 동적으로 생성됩니다 (그리고 트랜잭션 컨텍스트를 제공하는 데 필요한대로 `CacheKVStore`를 통해 래핑됩니다).

#### Tx, Msg, Handler 및 AnteHandler

트랜잭션 (`Tx` 인터페이스)은 서명/인증 된 메시지 (`Msg` 인터페이스)입니다.

Tendermint 메모리풀이 발견 한 트랜잭션은 트랜잭션의 유효성을 확인하고 수수료를 징수하는 AnteHandler (_ante_ 는 단지 이전(before)을 의미)에 의해 처리됩니다.

블록에서 커밋 된 트랜잭션은 먼저 `AnteHandler`를 통해 처리되고 수수료가 공제된 후 트랜잭션이 유효하면 해당 Handler를 통해 처리됩니다.

두 경우 모두 트랜잭션 바이트를 먼저 구문 분석해야합니다. 기본 트랜잭션 파서는 Amino를 사용합니다. 대부분의 SDK 개발자는 `x/auth` 패키지 (그리고 `x/auth`에서 제공되는 해당 AnteHandler 구현)에 정의 된 표준 트랜잭션 구조를 사용하고자할 것입니다:

```go
 // StdTx는 Msg를 요금 및 서명으로 랩핑하는 표준 방법입니다.
 // 참고: 첫 번째 서명은 FeePayer (서명은 nil이 아니어야 함)입니다.
 type StdTx struct {
    Msg        sdk.Msg        `json:"msg"`
    Fee        StdFee         `json:"fee"`
    Signatures []StdSignature `json:"signatures"`
 }
```

다양한 패키지는 일반적으로 자체 메시지 유형을 정의합니다.  `Basecoin` 예제 애플리케이션은 `app.go`에 등록 된 여러 메시지 유형을 포함합니다:

```go
sdk.RegisterWire(cdc)    // Msgs 등록
 bank.RegisterWire(cdc)
 stake.RegisterWire(cdc)
 slashing.RegisterWire(cdc)
 ibc.RegisterWire(cdc)
```

마지막으로 핸들러는 `app.go` 파일의 라우터에 추가되어 메시지를 해당 핸들러에 매핑합니다. (앞으로 더 많은 유연성을 위해 패턴 매칭을 가능하게하는 더 많은 라우팅 기능을 제공 할 예정입니다.)

```go
 // 메시지 라우트 등록
 app.Router().
    AddRoute("auth", auth.NewHandler(app.accountMapper)).
    AddRoute("bank", bank.NewHandler(app.coinKeeper)).
    AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.coinKeeper)).
    AddRoute("stake", stake.NewHandler(app.stakeKeeper))
```

#### EndBlocker

`EndBlocker` 훅은 각 블록의 끝에 수행 될 콜백 로직을 등록 할 수 있게 합니다.  이렇게하면 유효성 검사기 인플레이션 Atom 조항 처리와 같은 백그라운드 이벤트를 처리 할 수 ​​있습니다:

```go
// 검증인 공급 처리
 blockTime := ctx.BlockHeader().Time // XXX 초로 가정하고 확인합니다.
 if pool.InflationLastTime+blockTime >= 3600 {
    pool.InflationLastTime = blockTime
    pool = k.processProvisions(ctx)
 }
```

그런데, SDK는 Cosmos 허브에 대한 모든 결합/비결합 기능을 제공하는 [스테이킹 모듈] (https://github.com/cosmos/cosmos-sdk/tree/develop/x/stake)을 제공합니다.

#### 작업 시작

시작하려면, 다음과 같은 간단한 단계를 따라야합니다:

1. [Cosmos-SDK] (https://github.com/cosmos/cosmos-sdk/tree/develop) 저장소 복제
2. 아직 존재하지 않는 애플리케이션에 필요한 모듈을 코딩하십시오.
3. 앱 디렉토리를 만듭니다. 앱 메인 파일에서, 필요한 모듈을 가져 와서 다른 KVStore들을 인스턴스화합니다.
4. 블록 체인을 실행하십시오.

파이처럼 쉽게! 소개가 끝나면 실습을 배우고 예제를 사용하여 SDK 애플리케이션을 코딩하는 방법을 배웁니다.

## 설정

### 전제 조건

- [go] (https://golang.org/dl/) 및 [git] (https://git-scm.com/downloads)이 설치되어 있어야 합니다.
- `PATH` 와 `GOPATH`를 설정하는 것을 잊지 마십시오.

### 작업 환경 설정

[Cosmos-SDK repo] (https://githum.com/cosmos/cosmos-sdk)로 가서 포크하십시오. 그런 다음 터미널을 열고:

```bash
cd $GOPATH/src/github.com/your_username
git clone https://github.com/your_username/cosmos-sdk.git
cd cosmos-sdk
```

이제 멋진 기능이나 모듈이 병합될 경우에 대비해 원본 Cosmos-SDK를 업스트림으로 추가합니다:

```bash
git remote add upstream https://github.com/cosmos/cosmos-sdk.git
git fetch upstream
git rebase upstream/master
```

또한 모듈 전용 브랜치를 만듭니다:

```bash
git checkout -b my_new_application
```

모두 설정했습니다!

## 앱 설계

간단한 관리 애플리케이션

이 자습서에서는 **간단한 관리 모듈**과 함께 **간단한 관리 응용프로그램**을 코딩 할 것입니다. Cosmos-SDK에서 작동하는 애플리케이션을 빌드하는데 필요한 기본 개념의 대부분이 설명될 것입니다. 참고로 이것은 Cosmos Hub에 사용되는 관리 모듈이 아닙니다. Hub에는 훨씬 더 많은 [고급 관리 모듈] (https://github.com/cosmos/cosmos-sdk/tree/develop/x/gov)이 대신 사용될 것입니다.

`simple_governance` 애플리케이션의 모든 코드는 [여기] (https://github.com/gamarin2/cosmos-sdk/tree/module_tutorial/examples/simpleGov/x/simple_governance)에서 찾을 수 있습니다. 모듈과 앱은 저장소의 루트 레벨이 아니라 examples 디렉토리에 있습니다. 이것은 단지 편의를 위한 것이며, 모듈과 애플리케이션을 루트 디렉토리에서 코딩 할 수 있습니다.

두말 할 것 없이, 같이 해보시죠!

### 요구사항

모듈의 요구사항을 적어 두는 것으로 시작하겠습니다. 우리는 다음과 같은 간단한 관리 모듈을 설계하고 있습니다:

- 코인 보유자가 각자 제출할 수있는 간단한 텍스트 제안서.
- 제안서는 Atom코인으로된 보증금과 함께 제출해야 합니다. 보증금이 `MinDeposit` 보다 큰 경우, 관련 제안서는 투표 기간에 들어갑니다. 그렇지 않으면 거부됩니다. 
- Bonded Atom 보유자는 1 Bonded Atom 으로 1 표를 제안서에 투표 할 수 있습니다.
- Bonded Atom 보유자는 투표 할 때 `예`, `아니요` 및  `기권`의 3 가지 옵션 중에서 선택할 수 있습니다.
- 투표 기간이 종료 된 후 `아니오` 투표 보다 `예` 투표가 많으면 제안서가 채택됩니다. 그렇지 않으면 거부됩니다.
- 투표 기간은 2 주입니다.

모듈을 설계할 때 특정 방법론을 채택하는 것이 좋습니다. 블록 체인 어플리케이션은 단지 복제된 상태 머신이라는 것을 기억하십시오. '상태'란 주어진 시간에 애플리케이션을 표현한 것입니다. 애플리케이션 개발자는 애플리케이션의 목표에 따라 '상태'가 무엇을 나타내는지 정의 할 수 있습니다. 예를 들어, 간단한 암호화폐 애플리케이션의 '상태'는 주소를 잔액에 매핑하는 것입니다.

미리 정의된 규칙에 따라 상태를 업데이트 할 수 있습니다. 상태와 트랜잭션이 주어지면 상태-머신 (즉, 애플리케이션) 이 새로운 상태를 반환합니다. 블록체인 애플리케이션에서 트랜잭션은 블록으로 묶이지만, 이론은 동일합니다. 상태와 트랜잭션 집합 (블록)이 주어지면 애플리케이션은 새로운 상태를 반환합니다. SDK 모듈은 애플리케이션의 일부일 뿐이지만, 동일한 원칙을 기반으로 합니다. 결과적으로, 모듈 개발자는 상태 전환을 유발하는 상태의 일부와 트랜잭션 타입의 일부를 정의 하기만 하면 됩니다.

요약하면 다음을 정의해야합니다:

- 애플리케이션의 현재 상태의 일부를 나타내는 `상태`.
- 상태 전환을 유발하는 메시지가 포함되는 `트랜잭션들`.

### 상태

여기서는 멀티 스토어의 스토어들뿐 아니라 필요한 타입 (트랜잭션 타입 제외)을 정의합니다.

우리의 투표 모듈은 매우 간단합니다. 단 하나의 유형, 즉 `제안`만 필요합니다. `제안들`은 투표 대상 항목입니다. 모든 사용자가 제출할 수 있습니다. 보증금은 제공되어져야 합니다.

```go
type Proposal struct {
    Title           string          // 제안서 제목
    Description     string          // 제안서 설명
    Submitter       sdk.Address     // 제출자의 주소. 제안서가 수락되면 보증금은 반환해야합니다.
    SubmitBlock     int64           // 제안서가 제출되는 블록. 또한 투표 기간이 시작되는 블록.
    State           string          // 상태는 "Open", "Accepted"또는 "Rejected"일 수 있습니다.

    YesVotes        int64           // 찬성 투표 수
    NoVotes         int64           // 반대 투표 수
    AbstainVotes    int64           // 기권 투표 수
}
```

스토어 측면에서는, 멀티 스토어에 `제안서`를 저장하기 위해 [KVStore] (# kvstore)를 하나만 만들 것입니다. 또한 각 제안서에 각 유권자가 선택한 `투표` (`찬성`, `반대` 또는 `거부`)를 저장합니다.


### 메시지

모듈 개발자로서, 정의해야 할 것은 `트랜잭션`이 아니라 `메시지`입니다. 트랜잭션과 메시지는 Cosmos-SDK에 존재하지만 트랜잭션은 메시지가 트랜잭션에 포함된다는 점에서 메시지와 다릅니다. 트랜잭션은 메시지를 감싸고 서명 및 수수료와 같은 일반적인 정보를 추가합니다. 모듈 개발자는, 트랜잭션은 걱정할 필요 없고 메시지만 신경쓰면 됩니다.

상태를 수정하기 위해 필요한 메시지를 정의하겠습니다. 위의 요구 사항에 따라 두 가지 유형의 메시지를 정의해야합니다: 

-`SubmitProposalMsg`: 제안서 제출
-`VoteMsg`: 제안서에 투표

```go
type SubmitProposalMsg struct {
    Title           string      // 제안서 제목
    Description     string      // 제안서 설명
    Deposit         sdk.Coins   // 보증금은 제출자가 지불합니다. 투표 기간에 들어가려면 MinDeposit 보다 커야합니다.
    Submitter       sdk.Address // 제출자의 주소
}
```

```go
type VoteMsg struct {
    ProposalID  int64           // 제안서의 ID입니다.
    Option      string          // 투표자가 선택한 옵션
    Voter       sdk.Address     // 투표자의 주소
}
```

## 구현

이제 타입을 정의 했으므로 실제로 애플리케이션을 구현할 수 있습니다.

SDK 포크의 루트에서 `app` 및 `cmd` 폴더를 만듭니다. 이 폴더에서, 애플리케이션의 메인 파일인 `app.go`와 애플리케이션의 REST 및 CLI 명령을 처리하는 저장소를 생성합니다. 

```bash
mkdir app cmd 
mkdir -p cmd/simplegovcli cmd/simplegovd
touch app/app.go cmd/simplegovcli/main.go cmd/simplegovd/main.go
```

튜토리얼의 뒷부분에서 이 파일들을 처리 할 것입니다. 첫 번째 단계는 간단한 관리 모듈을 처리하는 것입니다.

### 간단한 통제 모듈

먼저, 모듈 폴더로 가서 모듈을 위한 폴더를 만듭니다.

```bash
cd x/
mkdir simple_governance
cd simple_governance
mkdir -p client/cli client/rest
touch client/cli/simple_governance.go client/rest/simple_governance.go errors.go handler.go handler_test.go keeper_keys.go keeper_test.go keeper.go test_common.go test_types.go types.go wire.go
```

우리가 필요로 하는 파일을 추가하는 것으로 시작합시다. 모듈의 폴더는 다음과 유사해야합니다:

```
x
└─── simple_governance
      ├─── client
      │     ├───  cli
      │     │     └─── simple_governance.go
      │     └─── rest
      │           └─── simple_governance.go
      ├─── errors.go
      ├─── handler.go
      ├─── keeper_keys.go
      ├─── keeper.go
      ├─── types.go
      └─── wire.go
```

이 파일들 각각에 대해 자세히 살펴봅시다.

#### 타입 (`types.go`)

이 파일에서는, 모듈의 사용자 정의 타입들을 정의합니다. 여기에는 [상태](#State) 섹션의 유형과 [트랜잭션](#Transactions) 섹션에 정의된 사용자 정의 메시지 유형이 포함됩니다.

메시지가 아닌 새로운 유형의 경우, 애플리케이션의 컨텍스트에서 의미가있는 메소드를 추가 할 수 있습니다. 우리의 경우, 우리는 투표 메시지가 들어올 때 주어진 제안서의 집계를 쉽게 업데이트 할 수 있는 `updateTally` 함수를 구현할 것입니다.

메시지는 약간 다릅니다. 그들은 SDK의 `types` 폴더에 정의된 `Message` 인터페이스를 구현합니다. 
