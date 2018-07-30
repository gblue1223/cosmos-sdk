# 튜토리얼: Cosmos-SDK 어플리케이션을 코딩하는 방법

이 튜토리얼에서는 Cosmos-SDK를 사용하여 어플리케이션을 코딩하는 기본 사항을 학습합니다. 텐더민트 (Tendermint)와 코스모스 생태계 (Cosmos Ecosystem)에 대한 소개부터 시작하여 텐더민트 (Tendermint) 소프트웨어와 Cosmos-SDK 프레임워크에 대한 개괄적인 개요를 살펴 보겠습니다. 그 다음에는 코드에 대해 알아보고 Cosmos-SDK 제일 윗단에서 간단한 어플리케이션을 작성하는 방법을 안내합니다.

## 텐더민트와 코스모스

블록체인은 세 가지 개념적 레이어로 나눌 수 있습니다:

- **네트워킹:** 트랜잭션 전파에 책임이 있습니다.
- **합의:** 다음 트랜잭션 집합의 처리에 동의하기 위해 검증인 노드들을 활성화 합니다 (예: 블록체인에 트랜잭션 블록 추가).
- **어플리케이션:** 트랜잭션 집합에 주어진 상태를 업데이트하는 작업을 담당합니다. 즉, 트랜잭션을 처리합니다.

*네트워킹* 계층은 각 노드가 트랜잭션을 수신 하는지 확인합니다. *합의* 계층은 각 노드가 자신의 로컬 상태를 수정하기 위해 동일한 트랜잭션에 동의하는지 확인합니다. *어플리케이션* 레이어는 트랜잭션을 처리합니다. 트랜잭션과 상태가 주어지면 어플리케이션은 새로운 상태를 반환합니다. Bitcoin을 예를 들면, 어플리케이션 상태는 각 계정에 대한 원장 또는 잔액 목록입니다 (실제로는 UTXO의 목록으로 Unspent Transaction Output의 약자이지만 단순화를 위해 잔액이라고합시다). 그리고 트랜잭션은 이러한 잔액 목록을 변경하여 어플리케이션의 상태를 수정합니다. Ethereum의 경우 어플리케이션은 가상 머신입니다. 각 트랜잭션은 이 가상 머신을 통해 진행되며 어플리케이션의 상태가 트랜잭션 내부에서 호출된 스마트 계약에 따라 수정됩니다. 

텐더민트 이전에는, 블록체인을 만들려면 세 개의 레이어를 처음부터 모두 만들어야했습니다. 대부분의 개발자가 Bitcoin 코드베이스를 포크하거나 복제하는 것을 선호했지만 Bitcoin의 프로토콜의 한계에 의해 제약을받기 때문에 실증나는 작업이었습니다. Ethereum Virtual Machine (EVM)은 이 문제를 해결하고 스마트 계약을 통해 사용자 정의 가능한 논리를 실행할 수 있게 하여 분산된 어플리케이션 개발을 단순화하도록 설계되었습니다. 그러나 블록체인 자체의 한계 (상호 운용성, 확장 성 및 사용자 정의)를 해결하지 못했습니다. Go-Ethereum은 Bitcoin의 코드베이스처럼 하드 포크하기가 어려운 매우 획일적인 기술 스택입니다. 

텐더민트는 이러한 문제를 해결하고 개발자에게 대안을 제공하기 위해 고안되었습니다. 텐더민트의 목표는 빌드 하려는 모든 어플리케이션 개발자를 지원하는 일반 엔진으로서 블록체인의 *네트워킹* 및 *합의* 레이어를 제공하는 것입니다. 텐더민트를 사용하면 개발자는 *어플리케이션* 레이어에만 집중해야하므로 수백시간의 작업시간과 값 비싼 개발 환경을 절약 할 수 있습니다. 참고로 텐더민트는 텐더민트 코어 엔진 내에서 사용 되는 비잔틴 장애 허용 합의 알고리즘을 나타냅니다.

텐더민트는, [ABCI](https://github.com/tendermint/abci)(Application-BlockChain-Interface) 라 불리는 소켓 프로토콜을 통해 블록체인 엔진인 텐더민트 코어 (*네트워킹* 및 *합의* 레이어)를 *어플리케이션* 레이어에 연결합니다. 개발자는 텐더민트 코어 엔진에서 실행되는 ABCI 지원 어플리케이션을 작성하기 위해 몇 가지 메시지만 구현 하면 됩니다. ABCI는 언어에 구속력이 없습니다. 즉, 개발자는 어떤 프로그래밍 언어를 이용해서 블록체인의 *어플리케이션* 파트를 구현할 수 있습니다. 텐더민트 코어 엔진 위에 구축하면 다음과 같은 이점도 있습니다.

- **공용 또는 사설 블록체인이 가능합니다.** 개발자는 텐더민트 코어 위에 권한이 부여 된 (사설) 권한이 없는 (공용) 모든 블록체인 어플리케이션을 배포 할 수 있습니다. 
- **성능.** 텐더민트 코어는 짧은 시간 간격으로 많은 수의 트랜잭션을 처리 할 수있는 최첨단 블록체인 컨센서스 엔진입니다. 텐더민트 코어의 블록 시간은 1 초 정도로 낮을 수 있으며 해당 기간에 수천 회의 트랜잭션을 처리 할 수 ​​있습니다.
- **즉각적인 완결성.** 텐더민트 컨센서스 알고리즘의 속성은 즉각적인 완결성입니다. 즉 유효성 검사기의 1/3 미만이 악성 코드 (비잔틴) 인 한 포크는 생성되지 않습니다. 사용자는 블록이 생성되는 즉시 트랜잭션이 완결되었는지 확인할 수 있습니다.
- **보안.** 텐더민트 코어의 합의는 단지 장애를 허용하는 것이 아니라, 책임감있는 최적의 비잔틴 장애 허용 (BFT) 에 있습니다. 만일 블록체인이 포크되면, 책임을 밝힐 수있는 방법이 있습니다.
- **라이트 클라이언트 지원.** 텐더민트는 내장된 라이트 클라이언트를 제공합니다.

하지만 가장 중요한 점은 텐더민트는 [Inter-Blockchain Communication Protocol](https://github.com/cosmos/cosmos-sdk/tree/develop/docs/spec/ibc) (IBC)과 호환된다는 것입니다. 즉, 공용이든 사설이든 텐더민트 기반 블록체인은 본질적으로 코스모스 생태계와 연결되어 생태계의 다른 블록체인과 토큰을 안전하게 교환 할 수 있습니다. IBC와 코스모스를 통한 상호 운용성의 혜택은 텐더민트 체인의 자주권을 보호합니다. 비 텐더민트 체인은 IBC 어댑터 또는 페그존을 통해 코스모스에 연결할 수도 있지만 이 내용은 이 문서의 범위를 벗어납니다.

코스모스 생태계에 대한 자세한 내용은 [이 글](https://blog.cosmos.network/understanding-the-value-proposition-of-cosmos-ecaef63350d)을 참조하십시오.

## 코스모스SDK 소개

텐더민트 기반 블록체인을 개발한다는 것은 어플리케이션 (즉, 상태 머신) 만 코딩 하면 된다는 것을 의미합니다. 하지만 그 자체로는 다소 어려울 수 있습니다. 이것이 바로 Cosmos-SDK가 존재하는 이유입니다.

[Cosmos-SDK](https://github.com/cosmos/cosmos-sdk)는 Cosmos 허브와 같은 다중 자산 Proof-of-Stake (PoS) 블록체인이자 Proof-Of -Authority (PoA) 블록체인을 구축하기위한 플랫폼.

Cosmos-SDK의 목표는 개발자가 일반적인 블록체인 기능을 다시 만들 필요 없이 Cosmos 네트워크 내에서 상호 운용 가능한 맞춤 블록 첸인 어플리케이션을 쉽게 만들 수있게하고, 텐더민트 ABCI 어플리케이션을 작성하는 복잡성을 제거하는 것입니다. 우리는 텐더민트 위에 안전한 블록체인 어플리케이션을 구축하기 위해 npm과 유사한 프레임워크의 SDK를 구상하고 있습니다.

SDK는 설계 측면에서 유연성과 보안성을 최대한 신경쓰고 있습니다. 프레임워크는 어플리케이션이 원하는대로 요소를 혼합하고 일치시킬 수있는 모듈 실행 스택 위주로 설계되었습니다. 또한 모든 모듈은 보다 강력한 어플리케이션 보안을 위해 샌드박스화되어 있습니다.

이것은 두 가지 주요 원칙에 기반합니다:

- **합성성(Composability):** 누구나 Cosmos-SDK 용 모듈을 만들 수 있으며 이미 구축 된 모듈을 통합하는 것은 블록체인 어플리케이션으로 가져 오는 것 만큼 간단합니다.

- **기능들Capabilities):** SDK는 기능 기반 보안에서 영감을 얻었으며, 수년간 블록체인 상태 머신과의 씨름을 통해 영향을 끼쳤습니다. 대부분의 개발자는 자체 모듈을 만들 때 다른 타사 모듈에 액세스 해야합니다. Cosmos-SDK는 개방형 프레임워크이며 이러한 모듈 중 일부는 악의적이라고 가정하기 때문에, 객체 기능 (OCAPS) 기반 원칙을 사용하여 SDK를 설계했습니다. 실제로, 이는 각 모듈이 다른 모듈에 대한 액세스 제어 목록을 유지하는 대신, 각 모듈은 사전 정의 된 기능 세트를 부여하기 위해 다른 모듈에 전달할 수 있는 키퍼라는 특수 객체를 구현한다는 것을 의미합니다. 예를 들어 모듈 A의 키퍼 인스턴스가 모듈 B로 전달되면, 후자는 모듈 A의 제한된 함수 집합을 호출 할 수 있습니다. 각 키퍼의 기능은 모듈 개발자가 정의하며 각 타사 모듈에 전달되는 기능을 기반으로 타사 모듈의 외부 코드 안전성을 이해하고 감사하는 것은 개발자의 임무입니다. '기능'에 대한 더 자세한 내용은 [이 글](http://habitatchronicles.com/2017/05/what-are-capabilities/)을 참조하십시오.

*참고: 현재 Cosmos-SDK는 Golang에만 존재합니다. 즉, 개발자는 Golang으로만 SDK 모듈을 개발할 수 있습니다. 앞으로, SDK는 다른 프로그래밍 언어로 구현 될 것이라 기대됩니다. 텐더민트팀에서 자금을 지원하는 것도 결국에는 가능할 것입니다.*

## 텐더민트 및 ABCI에 대한 알림

Cosmos-SDK는 블록체인의 *어플리케이션* 레이어를 개발하기위한 프레임워크입니다. 이 어플리케이션은 [Application-Blockchain Interface](https://github.com/tendermint/abci)의 약자 인 ABCI라는 간단한 프로토콜을 지원하는 모든 합의 엔진 (*합의* + *네트워킹* 레이어)에 연결될 수 있습니다.

텐더민트 코어는 Cosmos-SDK 위에 구축된 기본 컨센서스 엔진입니다. *어플리케이션* 및 *합의 엔진* 각자의 책임을 잘 이해하는 것이 중요합니다.

*합의 엔진*의 책임:
- 트랜잭션 전파
- 유효한 트랜잭션의 순서에 대한 동의.

*어플리케이션*의 책임:
- 트랜잭션 생성
- 트랜잭션이 유효한지 확인.
- 트랜잭션 처리 (상태 변경 포함)

*합의 엔진*은 각 블록에 대해 제공된 유효성 검증인 집합에 대한 지식을 갖고 있지만 유효성 검사기 집합 변경을 발동하는 것은 *어플리케이션*의 책임이라는 점을 강조 할 가치가 있습니다. 이것이 Cosmos-SDK와 텐더민트를 통해 **공개 체인과 사설 체인**을 모두 구축 할 수 있는 이유입니다. 체인은 유효성 검증인의 설정 변경을 제어하는 ​​어플리케이션 수준에서 정의 된 규칙에 따라 공개 또는 사설이됩니다.

ABCI는 *합의 엔진*과 *어플리케이션* 사이의 연결을 설정합니다. 본질적으로, 핵심은 두 가지 메시지입니다:

- `CheckTx`: 트랜잭션이 유효한지 어플리케이션에게 물어 봅니다. 검증인의 노드가 트랜잭션을 받으면 검증인의 노드는 `CheckTx`를 실행합니다. 트랜잭션이 유효하면 메모리풀에 추가됩니다.
- `DeliverTx`: 어플리케이션에게 트랜잭션을 처리하고 상태를 업데이트하도록 요청합니다.

*합의 엔진*과 *어플리케이션*이 서로 어떻게 상호 작용 하는지에 대한 개괄적인 개요를 설명하겠습니다.

- 항상, 검증인 노드의 합의 엔진 (텐더민트 코어)이 트랜잭션을 수신 할 때마다 CheckTx를 통해 어플리케이션에 전달하여 유효성을 검사합니다. 유효하면 트랜잭션이 메모리풀에 추가됩니다.
- 우리가 블록 N에 있다고합시다. 유효성 검증인 세트 V가 있습니다. 다음 블록의 제안자는 *합의 엔진*에 의해 V에서 선택됩니다. 제안자는 메모리풀에서 유효한 트랜잭션을 수집하여 새 블록을 만듭니다. 그런 다음 블록은 다른 검증인들에게 알려지고 서명/커밋 됩니다. V의 2/3 이상이 *사전커밋*에 서명하면 블록은 블록 N+1이됩니다 (합의 알고리즘에 대한 자세한 설명은 [여기](https://github.com/tendermint/tendermint/wiki/Byzantine-Consensus-Algorithm)를 클릭하십시오.
- 블록 N+1이 V의 2/3 이상에 의해 서명 될 때, 그것은 full-node로 알려지게 됩니다. full-node가 해당 블록을 수신하면 유효성을 확인합니다. 블록이 V의 2/3 이상 유효한 서명들을 보유하고 블록의 모든 트랜잭션들이 유효한 경우 블록은 유효합니다. *합의 엔진*은 트랜잭션의 유효성을 검사하기 위해 'DeliverTx'를 통해 어플리케이션으로 전송합니다. 각 트랜잭션 이후에 'DeliverTx`는 트랜잭션이 유효하면 새로운 상태를 반환합니다. 블록이 끝나면 최종 상태가 확약됩니다. 물론, 이것은 블록 내의 트랜잭션 순서가 중요하다는 것을 의미합니다.

## SDK-app의 아키텍처

Cosmos-SDK는 텐더민트 기반 블록체인 어플리케이션을 위한 기본 템플릿을 제공합니다. 이 템플릿은 [여기](https://github.com/cosmos/cosmos-sdk)에서 찾을 수 있습니다.

본질적으로, 블록체인 어플리케이션은 단순히 복제된 상태 머신입니다. 상태 (예: 암호화폐, 각 계정이 보유하는 동전 수) 및 상태 전이를 트리거하는 트랜잭션이 있습니다. 어플리케이션 개발자는 상태, 트랜잭션 유형 및 다른 트랜잭션이 상태를 수정하는 방법을 정의합니다.

### 모듈성

Cosmos-SDK는 모듈 기반 프레임워크입니다. 각 모듈은 그 자체로 작은 상태 기계이며 다른 모듈과 쉽게 결합되어 일관된 어플리케이션을 생성 할 수 있습니다. 즉, 모듈은 전역 상태와 트랜잭션 유형의 하위 섹션을 정의합니다. 그런 다음, 각 유형에 따라 올바른 모듈로 트랜잭션을 라우팅하는 것은 루트 어플리케이션의 일입니다. 이 과정을 이해하기 위해 상태 머신의 단순화된 일반적인 주기를 살펴 보겠습니다.

텐더민트 코어 엔진에서 트랜잭션을 수신하면, 다음은 *어플리케이션*이 수행하는 작업입니다:

1. 트랜잭션을 디코드하고 메시지를 받음
2.`Msg.Type()` 메소드를 사용하여 메시지를 적절한 모듈로 라우팅
3. 모듈에서 트랜잭션을 실행. 트랜잭션이 유효하면 상태를 수정.
4. 새 상태 또는 오류 메시지를 반환

단계 1, 2 및 4는 루트 어플리케이션에서 처리합니다. 3 단계는 해당 모듈에서 처리합니다.

### SDK 구성 요소

이를 염두에두고 SDK의 중요한 디렉토리를 살펴 보겠습니다:

- `baseapp`: 기본 어플리케이션을위한 템플릿을 정의. 기본적으로 Cosmos-SDK 응용 프록로그램이 기본 텐더민트 노드와 통신 할 수 있도록 ABCI 프로토콜을 구현합니다.
- `client`: 어플리케이션과 상호 작용하는 커맨드라인 인터페이스
- `server`: 노드와 통신하기 위한 REST 서버
- `examples`: `baseapp`과 모듈을 기반으로 작동하는 어플리케이션을 만드는 법에 대한 예제가 들어 있음
- `store`: 멀티 스토어에 대한 코드를 포함. 멀티 스토어는 기본적으로 사용자의 상태입니다. 각 모듈은 멀티 스토어에서 원하는 수의 KVStore들을 생성 할 수 있습니다. 해당 `키퍼`를 사용하여 각 상점에 대한 액세스 권한을 올바르게 처리하도록 주의하십시오.
- `types`: SDK 기반의 모든 어플리케이션에 필요한 공통 유형.
- `x`: 모듈이 있는 곳. 이 디렉토리에는 이미 빌드된 모듈이 모두 있습니다. 이러한 모듈 중 하나를 사용하려면, 어플리케이션에서 적절하게 가져와야합니다. [App - 모든 것을 하나로 연결] 섹션에서 우리는 사용방법을 보게 될 것입니다.

### 입문자들을 위한 코드런

#### KVStore

KVStore는 SDK 어플리케이션의 기본 지속성 계층을 제공합니다.

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

계정에 하나, IBC에 하나 등과 같이, 여러 개의 KVStore를 어플리케이션에 마운트 할 수 있습니다.

```go
 app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keyStake, app.keySlashing)
```

요청이있는 경우, KVStore의 구현은 각 쿼리에 대해 Merkle 증명을 제공 할 책임이 있습니다.

```go
 func (st *iavlStore) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
```

스토어들은 지속성 레벨에서 트랜잭션을 제공 할 수 있도록 캐시래핑 될 수 있습니니다 (iterator에서도 잘 지원됩니다). 이 기능은 "AnteHandler"가 트랜잭션에 대한 모든 관련 비용을 공제 한 후 트랜잭션 처리를 위한 트랜잭션 격리 계층을 제공하는 데 사용됩니다. 캐시래핑은 블록체인을위한 가상 머신 또는 스크립팅 환경을 구현할 때 유용 할 수 있습니다.

#### go-amino

Cosmos-SDK는 Go 형식을 Protobuf3 호환 바이트로 직렬화 및 비직렬화 하기 위해 [go-amino](https://github.com/cosmos/cosmos-sdk/blob/96451b55fff107511a65bf930b81fb12bed133a1/examples/basecoin/app/app.go#L97-L111)를 광범위하게 사용합니다.

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

예를 들어, `Basecoin` 예제 어플리케이션은 _Ed25519_와 _Secp256k1_ 키에 대해 알고 있는데, 그 이유는 아래의 앱의 'codec'에 의해 등록 되었기 때문입니다.

```go
wire.RegisterCrypto(cdc) // 암호를 등록.
```

Go-Amino에 대한 자세한 내용은 https://github.com/tendermint/go-amino를 참조하십시오.

#### 키들, 키퍼들, 그리고 매퍼들

Cosmos-SDK는 전체 어플리케이션을 구성하기 위해 함께 포함될 수 있는 라이브러리의 생태계를 활성화할 수 있도록 설계되었습니다. 이 생태계를보다 안전하게 만들기 위해 최소-권한 원칙에 따라 디자인 패턴을 개발했습니다.

`매퍼들` 과 `키퍼들`은 컨텍스트를 통해 KVStore 들에 대한 액세스를 제공합니다. 이 둘의 유일한 차이점은 `매퍼`가 낮은 수준의 API를 제공하기 때문에 일반적으로 `키퍼`는 다른 `키퍼들` 및 `매퍼들`에 대한 참조를 보유 할 수 있지만 그 반대는 아닙니다.

`매퍼들` 과  `키퍼들`은 직접 스토어들을 참조하지 않습니다.  그들은 오직 _key_ (아래의 `sdk.StoreKey`) 만을 가지고 있습니다:

```go
type AccountMapper struct {

    // 컨텍스트에서 스토어에 액세스하는데 사용된 (노출되지 않은) 키입니다.
    key sdk.StoreKey

    // 계정 실제 타입.
    proto Account

    // 계정의 바이너리 인코딩/디코딩을 위한 와이어 코덱.
    cdc *wire.Codec
 }
```

이런 식으로, 당신은 `app.go` 파일에 모든 것을 연결하고 어떤 컴포넌트가 어떤 스토어들과 다른 컴포넌트에 접근 할 수 있는지를 볼 수 있습니다.

```go
// accountMapper를 정의합니다.
 app.accountMapper = auth.NewAccountMapper(
    cdc,
    app.keyAccount, // 대상 스토어
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

`매퍼들` 과 `키퍼들`은 앱 초기화시 스토어가 알려지지 않아 스토어들에 대한 직접적인 참조를 가질 수 없습니다.  스토어는 모든 커밋 된 트랜잭션 (ABCI `DeliverTx`를 통해) 및 메모리풀 확인 트랜잭션 (ABCI `CheckTx`를 통해)에 대해 동적으로 생성됩니다 (그리고 트랜잭션 컨텍스트를 제공하는 데 필요한대로 `CacheKVStore`를 통해 래핑됩니다).

#### Tx, Msg, Handler 및 AnteHandler

트랜잭션 (`Tx` 인터페이스)은 서명/인증 된 메시지 (`Msg` 인터페이스)입니다.

텐더민트 메모리풀이 발견 한 트랜잭션은 트랜잭션의 유효성을 확인하고 수수료를 징수하는 AnteHandler (_ante_ 는 단지 이전(before)을 의미)에 의해 처리됩니다.

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

다양한 패키지는 일반적으로 자체 메시지 유형을 정의합니다.  `Basecoin` 예제 어플리케이션은 `app.go`에 등록 된 여러 메시지 유형을 포함합니다:

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

그런데, SDK는 Cosmos 허브에 대한 모든 결합/비결합 기능을 제공하는 [staking 모듈](https://github.com/cosmos/cosmos-sdk/tree/develop/x/stake)을 제공합니다.

#### 작업 시작

시작하려면, 다음과 같은 간단한 단계를 따라야합니다:

1. [Cosmos-SDK](https://github.com/cosmos/cosmos-sdk/tree/develop) 저장소 복제
2. 아직 존재하지 않는 어플리케이션에 필요한 모듈을 코딩하십시오.
3. 앱 디렉토리를 만듭니다. 앱 메인 파일에서, 필요한 모듈을 가져 와서 다른 스토어들을 인스턴스화합니다.
4. 블록체인을 실행하십시오.

파이처럼 쉽게! 소개가 끝나면 실습을 배우고 예제를 사용하여 SDK 어플리케이션을 코딩하는 방법을 배웁니다.

## 설정

### 전제 조건

- [go](https://golang.org/dl/) 및 [git](https://git-scm.com/downloads)이 설치되어 있어야 합니다.
- `PATH` 와 `GOPATH`를 설정하는 것을 잊지 마십시오.

### 작업 환경 설정

[Cosmos-SDK repo](https://githum.com/cosmos/cosmos-sdk)로 가서 포크하십시오. 그런 다음 터미널을 열고:

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

간단한 관리 어플리케이션

이 자습서에서는 **간단한 관리 모듈**과 함께 **간단한 관리 어플리케이션**을 코딩 할 것입니다. Cosmos-SDK에서 작동하는 어플리케이션을 빌드하는데 필요한 기본 개념의 대부분이 설명될 것입니다. 참고로 이것은 Cosmos Hub에 사용되는 관리 모듈이 아닙니다. Hub에는 훨씬 더 많은 [고급 관리 모듈](https://github.com/cosmos/cosmos-sdk/tree/develop/x/gov)이 대신 사용될 것입니다.

`simple_governance` 어플리케이션의 모든 코드는 [여기](https://github.com/gamarin2/cosmos-sdk/tree/module_tutorial/examples/simpleGov/x/simple_governance)에서 찾을 수 있습니다. 모듈과 앱은 저장소의 루트 레벨이 아니라 examples 디렉토리에 있습니다. 이것은 단지 편의를 위한 것이며, 모듈과 어플리케이션을 루트 디렉토리에서 코딩 할 수 있습니다.

두말 할 것 없이, 같이 해보시죠!

### 요구사항

모듈의 요구사항을 적어 두는 것으로 시작하겠습니다. 우리는 다음과 같은 간단한 관리 모듈을 설계하고 있습니다:

- 코인 보유자가 각자 제출할 수있는 간단한 텍스트 제안서.
- 제안서는 Atom코인으로된 보증금과 함께 제출해야 합니다. 보증금이 `MinDeposit` 보다 큰 경우, 관련 제안서는 투표 기간에 들어갑니다. 그렇지 않으면 거부됩니다. 
- Bonded Atom 보유자는 1 Bonded Atom 으로 1 표를 제안서에 투표 할 수 있습니다.
- Bonded Atom 보유자는 투표 할 때 `예`, `아니요` 및  `기권`의 3 가지 옵션 중에서 선택할 수 있습니다.
- 투표 기간이 종료 된 후 `아니오` 투표 보다 `예` 투표가 많으면 제안서가 채택됩니다. 그렇지 않으면 거부됩니다.
- 투표 기간은 2 주입니다.

모듈을 설계할 때 특정 방법론을 채택하는 것이 좋습니다. 블록체인 어플리케이션은 단지 복제된 상태 머신이라는 것을 기억하십시오. '상태'란 주어진 시간에 어플리케이션을 표현한 것입니다. 어플리케이션 개발자는 어플리케이션의 목표에 따라 '상태'가 무엇을 나타내는지 정의 할 수 있습니다. 예를 들어, 간단한 암호화폐 어플리케이션의 '상태'는 주소를 잔액에 매핑하는 것입니다.

미리 정의된 규칙에 따라 상태를 업데이트 할 수 있습니다. 상태와 트랜잭션이 주어지면 상태-머신 (즉, 어플리케이션) 이 새로운 상태를 반환합니다. 블록체인 어플리케이션에서 트랜잭션은 블록으로 묶이지만, 이론은 동일합니다. 상태와 트랜잭션 집합 (블록)이 주어지면 어플리케이션은 새로운 상태를 반환합니다. SDK 모듈은 어플리케이션의 일부일 뿐이지만, 동일한 원칙을 기반으로 합니다. 결과적으로, 모듈 개발자는 상태 전환을 유발하는 상태의 일부와 트랜잭션 타입의 일부를 정의 하기만 하면 됩니다.

요약하면 다음을 정의해야합니다:

- 어플리케이션의 현재 상태의 일부를 나타내는 `상태`.
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

스토어 측면에서는, 멀티 스토어에 `제안서`를 저장하기 위해 [KVStore](#kvstore)를 하나만 만들 것입니다. 또한 각 제안서에 각 유권자가 선택한 `투표` (`찬성`, `반대` 또는 `거부`)를 저장합니다.


### 메시지

모듈 개발자로서, 정의해야 할 것은 `트랜잭션`이 아니라 `메시지`입니다. 트랜잭션과 메시지는 Cosmos-SDK에 존재하지만 트랜잭션은 메시지가 트랜잭션에 포함된다는 점에서 메시지와 다릅니다. 트랜잭션은 메시지를 감싸고 서명 및 수수료와 같은 일반적인 정보를 추가합니다. 모듈 개발자는, 트랜잭션은 걱정할 필요 없고 메시지만 신경쓰면 됩니다.

상태를 수정하기 위해 필요한 메시지를 정의하겠습니다. 위의 요구 사항에 따라 두 가지 유형의 메시지를 정의해야합니다: 

- `SubmitProposalMsg`: 제안서 제출
- `VoteMsg`: 제안서에 투표

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

이제 타입을 정의 했으므로 실제로 어플리케이션을 구현할 수 있습니다.

SDK 포크의 루트에서 `app` 및 `cmd` 폴더를 만듭니다. 이 폴더에서, 어플리케이션의 메인 파일인 `app.go`와 어플리케이션의 REST 및 CLI 명령을 처리하는 저장소를 생성합니다. 

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

메시지가 아닌 새로운 유형의 경우, 어플리케이션의 컨텍스트에서 의미가있는 메소드를 추가 할 수 있습니다. 우리의 경우, 우리는 투표 메시지가 들어올 때 주어진 제안서의 집계를 쉽게 업데이트 할 수 있는 `updateTally` 함수를 구현할 것입니다.

메시지는 약간 다릅니다. 그들은 SDK의 `types` 폴더에 정의된 `Message` 인터페이스를 구현합니다. 다음은 사용자 정의 메시지 타입을 정의할 때 구현해야하는 메소드입니다:

- `Type()`: 이 함수는 모듈의 경로 이름을 반환합니다. 메시지가 어플리케이션에 의해 처리될 때, 메시지들은 `Type()` 메쏘드에 의해 반환된 문자열을 사용하여 라우팅 됩니다.
- `GetSignBytes()`: 메시지의 바이트 표현을 반환합니다. 메시지에 서명하는데 사용됩니다.
- `GetSigners()`: 서명자(들)의 주소(들)을 반환합니다.
- `ValidateBasic()`: 이 함수는 명백히 유효하지 않은 메시지를 버리는데 사용됩니다. 이것은 baseapp 파일에서 `runTx()`의 시작부분에서 호출됩니다. `ValidateBasic()`이 `nil`을 반환하지 않으면 앱은 트랜잭션 실행을 중단합니다.
- `Get()`: 기본 getter는, 메시지의 속성을 반환합니다.
- `String()`: 사람이 읽을 수 있는 버전의 메시지를 반환합니다.

우리의 단순한 거버넌스 메시지의 의미는 다음과 같습니다:

- `Type()`은 `"simpleGov "`를 반환합니다.
- `SubmitProposalMsg`의 경우, 속성이 비어 있지 않고 보증금이 유효하고 양수인지 확인해야합니다. 이것은 기본적인 검증일 뿐이므로, 이 방법으로 송금인이 입금액을 지불 할 수 있는 충분한 자금을 가지고 있는지 확인하지 않습니다.
- `VoteMsg`의 경우, 주소와 옵션이 유효하고 proposal ID가 음수가 아닌지 확인합니다.
- 다른 방법들은, 사용자 정의가 필요하지 않습니다. 코드를 확인하여 표준 구현 방법을 확인할 수 있습니다.

#### 키퍼 (`keeper.go`)

##### 키퍼들에 대한 짧은 소개

`키퍼들`은 모듈 저장소에 대한 읽기/쓰기를 처리하는 모듈 추상화입니다. 이것은 Cosmos의 [`Object Capability Model`](링크)의 실용적인 구현입니다. 


모듈 개발자는, 우리의 모듈뿐만 아니라 다른 모듈을 위해서도 모듈 저장소와 상호 작용하도록 키퍼들을 정의해야합니다. 다른 모듈이 우리의 모듈 저장소 중 하나에 액세스 하려고 할때, 이 저장소의 키퍼를 어플리케이션 레벨에서 전달해야합니다. 실제적으로, 다음과 같이 보일 것입니다:

```go
// app.go 안에서

// 키퍼들을 인스턴스화 합니다.
keeperA = moduleA.newKeeper(app.moduleAStoreKey)
keeperB = moduleB.newKeeper(app.moduleBStoreKey)

// keeperA의 인스턴스를 모듈 B의 핸들러에 전달합니다.
app.Router().
        AddRoute("moduleA", moduleA.NewHandler(keeperA)).
        AddRoute("moduleB", moduleB.NewHandler(keeperB, keeperA)) // 여기서 모듈 B는 keeperA 인스턴스를 통해 모듈 A의 저장소 중 하나에 액세스 할 수 있습니다
```

`KeeperA`는 모듈 B의 처리기에 일련의 기능들을 부여합니다. 모듈을 개발할 때는, 키퍼들을 통해 부여 될 수 있는 다양한 기능들의 민감도를 생각하는 것이 좋습니다. 예를 들어, 일부 모듈은 모듈 A의 기본 저장소에 읽기와 쓰기가 필요한 반면 다른 모듈은 읽는 것만 필요할 수 있습니다. 모듈에 여러 상점이있는 경우 일부 키퍼들은 모든 상점에 대한 액세스 권한을 부여 할 수 있고, 다른 키퍼들은 특정 하위 상점에 대한 액세스 권한 만 부여 할 수 있습니다. 모듈 개발자는 어플리케이션 개발자가 올바른 기능들과 함께 키퍼를 쉽게 인스턴스화 할 수 있도록 확실히 해야합니다. 물론, 모듈의 핸들러는 그 모듈의 키퍼의 무제한적인 인스턴스를 얻을 가능성이 높습니다.

##### 우리 앱을 위한 스토어

키퍼자체에 대해 알아보기 전에, 거버넌스 하위 저장소에 어떤 객체들을 저장해야하고, 어떻게 색인을 만드는지 살펴 보겠습니다.

- `제안서`는 `'제안서'|<제안ID>`에 의해 색인됩니다.
- `'제안서'|<제안서ID>|'투표'|<투표자주소>`에 의해 `투표` (`찬성`,`반대`,`무효`)가 색인 됩니다.

`'제안서'` 와 `'투표'` 의 따옴표를 주목하십시오. 따옴표들은 이것들이 일정한 키워드라는 것을 나타냅니다. 예를 들어, 제안서 `0101`에 주소 `0x01`로 투표자가 캐스팅 한 옵션은 `'제안서'|0101|'투표'|0x01` 색인에 저장됩니다.

이러한 키워드는 범위 쿼리를 지원하는 데 사용됩니다. 범위 쿼리 (TODO: 공식 링크)를 통해 개발자는 상점의 부분 공간을 쿼리하고 반복자를 반환 할 수 있습니다. 그들은 백그라운드에서 사용되는 [IAVL+tree](https://github.com/tendermint/iavl)의 멋진 속성에 의해 가능합니다. 실제로, 이는 키-값 쌍의 지정된 부분 공간을 반복 하는 동안, O(1)의 속도로 키 - 값 쌍을 저장하고 쿼리 할 수 ​​있음을 의미합니다. 예를 들어, 우리는 `rangeQuery(SimpleGovStore, <제안서ID|'주소'>)`를 호출하여 주어진 제안서에 투표한 모든 주소를, 투표와 함께 쿼리 할 수 ​​있습니다.

##### 앱을 위한 키퍼들

우리의 경우, 우리는 단지 하나의 저장소인 `SimpleGov` 저장소만 가지고 있습니다. 우리는 우리의 키퍼들을 통해 이 상점 내부에 값을 설정하고 가져와야 할 것입니다. 그러나, 이 두 가지 조치는 보안 측면에서 동일한 영향을 미치지 않습니다. 저장소에 대한 읽기 권한을 다른 모듈에 부여하는 데는 아무런 문제가 없지만, 쓰기 권한은 훨씬 민감합니다. 이상적으로 어플리케이션 개발자는 저장소에서 값을 가져오기만 할 수 있는 관리 매퍼와, 값을 가져 오거나 설정할 수 있는 관리 매퍼를 만들 수 있어야합니다. 이를 위해 우리는 두 개의 키퍼들을 소개 할 것입니다: `Keeper` 와 `KeeperRead`. 어플리케이션 개발자가 어플리케이션을 만들면, 개발자들은 어떤 모듈 키키퍼를 사용할지 결정할 수 있습니다.

이제 우리 모듈의 키퍼가 어떤 **외부** 모듈의 키퍼를 액세스 해야 하는지에 대해 생각해 봅시다.
각 제안서에는 보증금이 필요합니다. 이 말은 우리 모듈이 토큰을 처리하는 모듈인 `bank` 모듈을 읽고 쓸 수 있어야 한다는 것을 의미합니다. 우리는 또한 각 지분에 근거하여 투표자의 투표권을 결정할 수 있어야 합니다. 이를 위해, `staking`모듈의 저장소에 대한 읽기 액세스가 필요합니다. 그러나, 이 저장소에 대한 쓰기 액세스는 필요하지 않습니다. 따라서 우리 모듈에서, 그리고 어플리케이션 개발자는 `staking` 모듈의 읽기 전용 키퍼만 우리 모듈의 핸들러에 전달하는지 주의해야 합니다.

모든 것을 염두에 두고, 우리는 `Keeper`의 구조를 정의 할 수 있습니다:

```go
    type Keeper struct {
        SimpleGov    sdk.StoreKey        // 모듈 저장소의 키
        cdc                 *wire.Codec         // 엔코드/디코드 구조체를 위한 코덱
        ck                  bank.Keeper         // 보증금 처리를 위해 필요. 이 모듈은 Atom 잔액에 대한 읽기/쓰기만 요청합니다.
        sm                  stake.Keeper        // 투표권 계산에 필요. 이 모듈은 staking 저장소에 대한 읽기 액세스만 필요합니다.
        codespace           sdk.CodespaceType   // 오류 코드용 공간을 예약
    }
```

그리고 우리의 `KeeperRead` 구조:

```go
type KeeperRead struct {
    Keeper
}
```

`KeeperRead`는 우리가 오버라이드한 것을 제외하고, `Keeper`로 부터 모든 메소드를 상속받습니다. 이것들은 상점에 쓰기를 수행하는 메소드들이 될 것입니다.

##### 함수 및 메서드

생성해야 하는 첫 번째 함수는 생성자입니다.

```go
func NewKeeper(SimpleGov sdk.StoreKey, ck bank.Keeper, sm stake.Keeper, codespace sdk.CodespaceType) Keeper
```

이 함수는 메인 `app.go` 파일에서 호출되어 새로운 `Keeper '를 생성합니다. 유사한 함수가 `KeeperRead`를 위해 종료됩니다.

```go
func NewKeeperRead(SimpleGov sdk.StoreKey, ck bank.Keeper, sm stake.Keeper, codespace sdk.CodespaceType) KeeperRead
```

어플리케이션 및 해당 모듈의 요구에 따라,  `Keeper`, `KeeperRead` 또는 둘 다 어플리케이션 수준에서 인스턴스화됩니다.

*주의: `Keeper` 타입 이름과 `NewKeeper()` 함수 이름은 모든 모듈에서 사용되는 표준 이름입니다. 이 표준을 따르지 않아도 되지만, 그렇게 하는게 어플리케이션 개발자로서서의 인생이 편해질 수 있습니다*

이제, 우리 모듈의 `Keeper`에 필요한 메소드를 설명하겠습니다. 완전한 구현을 확인하시려면 `keeper.go`를 참조하십시오.

- `GetProposal`: `제안서ID`가 주어진 `제안서`를 얻습니다. 제안서는 읽을 수 있으려면 `바이트`를 해독해야합니다.
- `SetProposal`: `'제인서들'|<제안서ID>` 색인에 `제안서`를 설정하십시오. 제안서는 저장되기 전에 `바이트`로 암호화되어야합니다.
- `NewProposalID`: 새로운 고유한 `proposalID`를 생성하는 함수입니다.
- `GetVote`: `proposalID` 와 `voterAddress`를 통해 투표 `Option`을 가가져옵니다.
- `SetVote`: `proposalID` 와 `voterAddress`를 통해 투표 `Option`을 설정합니다.
- 제안서 큐 메소드들: 이 메소드들은 FIFO(First-In First-Out) 기반으로 `제안서들`을 저장하기 위한 일반적인 제안서 큐를 구현합니다. 이는 투표 기간이 종료될 때 투표 수를 집계하는데 사용됩니다.

마지막으로 해야할 일은 `KeeperRead` 타입을 위한 특정 메소드를 오버라이드하는 것입니다. `KeeperRead`는 스토어에 대한 쓰기 권한이 없어야합니다. 그러므로 우리는 제안서 큐의 메소드들인 `SetProposal()`, `SetVote()` 및 `NewProposalID()` 메소드와 `setProposalQueue()` 메소드를 오버라이드 할 것입니다. `KeeperRead`의 경우, 이 메소드는 단지 에러를 던질 것입니다.

*주의: 코드를 보면, 컨텍스트 `ctx`가 많은 메소드의 매개 변수라는 것을 알 수 있습니다. 컨텍스트 `ctx '는 현재 블록 높이와 같은 현재 상태에 대한 유용한 정보를 제공하고 키퍼 `k`가 `KVStore`에 접근하도록 허용합니다. [여기](https://github.com/cosmos/cosmos-sdk/blob/develop/types/context.go#L144-L168) 에서 `ctx`의 모든 메소드를 확인할 수 있습니다*.

#### 핸들러 (`handler.go`)

##### 생성자와 핵심 핸들러들

핸들러들은 상태 머신의 핵심 로직을 구현합니다. 트랜잭션이 앱에서 모듈로 라우팅되면 `handler` 함수에 의해 실행됩니다.

실제로, 하나의 `handler '가 모듈의 각 메시지에 대해 구현 될 것입니다. 우리의 경우, 두 가지 메시지 타입이 있습니다. 그러므로 우리는 두 개의 `handler` 함수가 필요합니다. 우리는 또한 메시지를 올바른 `handler`에 라우팅하기 위해 생성자 함수가 필요합니다:

```go
func NewHandler(k Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        switch msg := msg.(type) {
        case SubmitProposalMsg:
            return handleSubmitProposalMsg(ctx, k, msg)
        case VoteMsg:
            return handleVoteMsg(ctx, k, msg)
        default:
            errMsg := "Unrecognized gov Msg type: " + reflect.TypeOf(msg).Name()
            return sdk.ErrUnknownRequest(errMsg).Result()
        }
    }
}
```

메시지는 타입에 따라 적절한 `핸들러`로 라우팅됩니다. 우리의 간단한 관리 모듈의 경우 두 개의 메시지 유형에 해당하는 두 개의 `핸들러`만 있습니다. 핸들러들은 비슷한 특징이 있습니다:

```go
func handleSubmitProposalMsg(ctx sdk.Context, k Keeper, msg SubmitProposalMsg) sdk.Result
```

이 함수의 매개 변수를 살펴 보겠습니다:

- 스토어에 액세스하기 위한 컨텍스트 `ctx '.
- 키퍼 `k`는 핸들러에게 모듈 스토어 (우리의 경우 `SimpleGovernance`)를 포함한 다른 스토어들과, 키퍼 `k`가 액세스 권한을 부여받은 (우리의 경우 `stake`와 `bank`) 다른 모듈들의 모든 스토어들을 읽고 쓸 수 있습니다.
- `msg` 메시지는 트랜잭션의 송신자가 제공 한 모든 정보를 보유하고 있습니다.

이 함수는 어플리케이션에 'Result'를 반환합니다. 여기에는 이 거래에 대한 `Gas` 금액 및 메시지가 성공적으로 처리 되었는지 여부와 같은 유용한 정보가 포함되어 있습니다. 이 시점에서, 우리는 간단 관리 모듈의 영역을 빠져 나와 루트 어플리케이션 레벨로 돌아갑니다. `Result`는 어플리케이션마다 다릅니다. 더 많은 정보를 원하시면 `sdk.Result` 타입을 [여기](https://github.com/cosmos/cosmos-sdk/blob/develop/types/result.go)에서 직접 확인하십시오.

##### BeginBlocker 및 EndBlocker

대부분의 스마트 계약 플랫폼과는 달리, Cosmos-SDK 어플리케이션에서 로직 실행을 자동으로 수행 할 수 있습니다 (즉, 최종 사용자가 전송 한 트랜잭션에 의해 트리거되지 않음).

이 자동 실행 코드는 모든 블록의 처음과 끝에 호출 되는 `BeginBlock`과`EndBlock` 함수에서 발생합니다. 그것들은 강력한 도구이지만, 어플리케이션 개발자는 이들을 조심해야합니다. 예를 들어, 비싼 계산이 블록 시간을 지연시킬 수 있고, 끝이 없는 루프가 체인을 완전히 동결시킬 수 있으므로, 개발자가 이 함수들에서 발생하는 컴퓨팅의 양을 제어하는 ​​것은 매우 중요한 일입니다.

`BeginBlock` 과 `EndBlock`은 각각의 모듈이 자신의 `BeginBlock`과`EndBlock` 로직을 구현할 수 있다는 것을 의미하는, 조합 가능한 함수입니다. 필요할 때 `BeginBlock` 과 `EndBlock` 로직이 모듈의 `handler`에서 구현됩니다. 다음은 `EndBlock`을 위한 표준적인 방법입니다 (`BeginBlock`도 똑같은 패턴을 따릅니다):

```go
func NewEndBlocker(k Keeper) sdk.EndBlocker {
    return func(ctx sdk.Context, req abci.RequestEndBlock) (res abci.ResponseEndBlock) {
        err := checkProposal(ctx, k)
        if err != nil {
            panic(err)
        }
        return
    }
}
```

각 모듈이 어플리케이션 레벨에서`BeginBlock`과`EndBlock` 생성자를 선언해야한다는 것을 잊지 마십시오. [어플리케이션 - 모두 함께 연결하기](# application _-_ bridging_it_all_together)를 참조하십시오.

우리의 단순한 관리 어플리케이션을 위해, 우리는 투표 결과를 자동으로 집계하기 위해 `EndBlock`을 사용할 것입니다. 수행할 단계는 다음과 같습니다:

1. `ProposalProcessingQueue`에서 가장 오래된 제안서를 받으십시오.
2. `CurrentBlock`이 제안에 대한 투표 기간이 끝나는 블록인지 확인합니다. 만일 그렇다면, 3으로 이동하십시오. 아니라면, 종료하십시오.
3. 제안서가 수락 또는 거부되었는지 확인하십시오. 제안서 상태를 업데이트하십시오.
4. 'ProposalProcessingQueue`에서 제안서를 꺼내고 1번으로 돌아갑니다.

이 프로세스에 대해 간단한 안전성 분석을 수행해봅시다.
- `ProposalProcessingQueue`에서 제안서의 수가 유한하기 때문에 루프가 영원히 돌아 가지 않습니다.
- 각 제안서의 집계가 비싸지 않고 제안의 수가 상대적으로 낮을 것으로 예상되므로 계산이 너무 비싸지 않아야합니다. 이는 제안서가 수락되려면 `보증금`이 필요하기 때문입니다. `MinDeposit`은 우리가 너무 많은 `Proposals`을 대기열에 넣지 않을 만큼 충분히 높아야합니다.
- 어플리케이션이 너무 성공적이어서 `ProposalProcessingQueue`가 블록체인이 느려지는 수많은 제안을 포함하게 되면, 모듈은 상황을 완화하기 위해 수정되어야합니다. 하나의 영리한 방법은 `MaxIteration`에서 개별 `EndBlock` 당 반복 횟수를 제한하는 것입니다. 이렇게 하면, 제안서의 수가 매우 중요하고 블록 시간이 안정적이어야하는 경우, 여러 블록에걸쳐 집계를 할 수 있습니다. 이것은 현재 검사 `if (CurrentBlock == Proposal.SubmitBlock + VotingPeriod)` 를 `if (CurrentBlock > Proposal.SubmitBlock + VotingPeriod) AND (Proposal.Status == ProposalStatusActive)`로 수정해야 할 것입니다.

#### 와이어 (`wire.go`)

`wire.go` 파일은 개발자가 모듈의 구체적인 메시지 유형을 코덱에 등록 할 수 있게합니다. 우리의 경우 선언할 두 가지 메시지가 있습니다:

```go
func RegisterWire(cdc *wire.Codec) {
    cdc.RegisterConcrete(SubmitProposalMsg{}, "simple_governance/SubmitProposalMsg", nil)
    cdc.RegisterConcrete(VoteMsg{}, "simple_governance/VoteMsg", nil)
}
```
`app.go` 에서 이 함수를 호출하는 것을 잊지 마십시오 (자세한 내용은 [Application - Bridging it all together](#application _-_ bridging_it_all_together) 참조).

#### 오류 (`errors.go`)

`error.go` 파일은 모듈에 대한 커스텀 에러 메시지를 정의 할 수 있게합니다.  오류를 선언하는 것은 모든 모듈에서 상대적으로 비슷해야합니다. 우리의 간단한 관리 모듈의 [error.go](./error.go) 파일에서 구체적인 예를 볼 수 있습니다. 코드자체가 설명을 대신해줄 것입니다.

우리 모듈의 에러는 `sdk.Error` 인터페이스를 상속하므로 `Result()` 메소드를 가집니다. 이 메소드는 `handler`에 에러가 있고 실제 결과 대신 에러가 리턴되어야 할 때 유용합니다.

#### 명령줄 인터페이스 및 Rest API

각 모듈은 명령줄 인터페이스용 명령 세트와 REST API 용 엔드포인트를 정의 할 수 있습니다. 

##### 명령줄 인터페이스 (CLI)

`cli` 폴더로 가서 `simple_governance.go` 파일을 만드십시오. 여기서 모듈에 대한 명령을 정의합니다.

CLI는 [Cobra](https://github.com/spf13/cobra) 위에 구축됩니다. 다음은 Cobra 위에 명령을 작성하는 스키마입니다:

```go
    // 플래그를 선언하십시오.
    const(
        Flag = "flag"
        ...
    )

    // 주요 명령 함수. 각 명령당 하나의 함수.
    func Command(codec *wire.Codec) *cobra.Command {
        // 반환할 명령을 만듭니다.
        command := &cobra.Command{
            Use: "actual command",
            Short: "Short description",
            Run: func(cmd *cobra.Command, args []string) error {
                // 명령이 사용될 때 실행할 실제 함수
            },
        }

        // 명령에 플래그를 추가
        command.Flags().<Type>(FlagNameConstant, <example_value>, "<Description>")

        return command
    }
```

간단한 관리 모듈의 명령들에 대한 상세 구현은 [여기](../client/cli/ simple_governance.go)를 클릭하십시오.

##### Rest API

Rest 서버 [Light-Client Daemon (LCD)](https://github.com/cosmos/cosmos-sdk/tree/master/client/lcd)는 **HTTP 쿼리**를 지원합니다.

________________________________________________________

사용자 인터페이스 <=======> REST SERVER <=======> FULL-NODE

________________________________________________________

완전한 노드를 실행하고 싶지 않은 최종 사용자가 체인과 교류할 수 있습니다. LCD는 `--trust-node` 플래그를 통해 **Light-Client 검증**을 수행하도록 구성할 수 있습니다. 이 옵션은 `true` 또는 `false`로 설정할 수 있습니다.

- *Light-Client 검증*을 사용하는 경우, Rest Server는 light-client로 작동하며 최종 사용자의 시스템에서 실행해야합니다. 이를 통해 체인 전체를 로컬에 저장하지 않고도 신뢰할 수 없는 방식으로 체인을 간섭 할 수 있습니다.

- *Light-Client 검증*이 비활성화 된 경우, Rest Server는 HTTP 호출에 대한 간단한 릴레이 역할을합니다. 이 설정에서, Rest Server는 최종 사용자의 컴퓨터에서 실행하지 않아도 됩니다. 대신, 서버가 연결되는 full-node를 운영하는 동일한 엔티티에 의해 실행됩니다. 이 모드는 최종 사용자가 full-node 운영자를 신뢰하고 로컬에 저장하지 않으려는 경우에 유용합니다.

이제 사용자가 HTTP 요청을 통해 쿼리할 수 있는 endpoint들을 정의해 보겠습니다. 이러한 endpoint들은 `rest` 폴더에 저장된 `simple_governance.go` 파일에 정의 될 것입니다.

| Method | URL                             | Description                                                 |
|--------|---------------------------------|-------------------------------------------------------------|
| GET    | /proposals                      | Range query to get all submitted proposals                  |
| POST   | /proposals                      | Submit a new proposal                                       |
| GET    | /proposals/{id}                 | Returns a proposal given its ID                             |
| GET    | /proposals/{id}/votes           | Range query to get all the votes casted on a given proposal |
| POST   | /proposals/{id}/votes           | Cast a vote on a given proposal                             |
| GET    | /proposals/{id}/votes/{address} | Returns the vote of a given address on a given proposal     |

프런트엔드 개발자와 서비스 공급자가 적절하게 상호 작용할 수 있도록 적절한 endpoints를  제공하는 것은 모듈 개발자의 임무입니다.

간단한 관리 모듈의 endpoints에 대한 실제 코드내 구현에 대해서는 [이 파일](../client/rest/simple_governance.go)을 살펴볼 수 있습니다. 또한 REST API 모범 사례를 위한 [링크](https://hackernoon.com/restful-api-designing-guidelines-the-best-practices-60e1d954e7c9)가 있습니다.

### 어플리케이션 - 모두 함께 연결하기

이제 우리가 필요한 모든 조각을 만들었습니다. 이제는 어플리케이션에 통합할 때입니다. `/ x`  디렉토리를 빠져 나와 SDK 디렉토리의 루트로 돌아가겠습니다.


```bash
// 디렉토리의 루트 레벨에서
cd app
```

이제 간단한 관리 어플리케이션을 만들 준비가되었습니다!

#### 어플리케이션 구조

* 참고 사항: 전체 파일을 확인할 수 있습니다 (댓글 포함). [이곳](link)*

`app.go` 파일은 여러분의 어플리케이션을 정의하는 메인 파일입니다. 그 안에는 필요한 모든 모듈, 키퍼, 핸들러, 스토어 등을 선언 할 것입니다. 이 파일의 각 섹션을 살펴보고 어플리케이션의 구성 방법을 확인하십시오.

둘째, 우리는 우리의 어플리케이션의 이름을 정의해야합니다.

```go
const (
    appName = "SimpleGovApp"
)
```

그런 다음 어플리케이션의 구조를 정의합시다.

```go
// 확장된 ABCI 어플리케이션
type SimpleGovApp struct {
    *bam.BaseApp
    cdc *wire.Codec

    // 서브스토어에 액세스하는 키
    capKeyMainStore      *sdk.KVStoreKey
    capKeyAccountStore   *sdk.KVStoreKey
    capKeyStakingStore   *sdk.KVStoreKey
    capKeySimpleGovStore *sdk.KVStoreKey

    // 키퍼
    feeCollectionKeeper auth.FeeCollectionKeeper
    coinKeeper          bank.Keeper
    stakeKeeper         simplestake.Keeper
    simpleGovKeeper     simpleGov.Keeper

    // 계정 가져 오기 및 설정 관리
    accountMapper auth.AccountMapper
}
```

- 각 어플리케이션은 `BaseApp` 템플리트, 즉 포인터 위에 빌드됩니다.
- `cdc`는 우리의 어플리케이션에서 사용된 코덱입니다.
- 그런 다음 어플리케이션에서 필요한 스토어의 키를 가져 오십시오. 우리의 간단한 관리 어플리케이션의 경우, 우리는 3 개의 스토어 + 메인 스토어가 필요합니다.
- 그런 다음, 키퍼와 매퍼를 불러옵니다.

우리가 왜 이러한 스토어들과 키퍼들이 필요한지 명확하게 하기 위해 간단히 상기시켜보겠습니다. 우리의 어플리케이션은 기본적으로 simple_governance 모듈을 기반으로합니다. 그러나 우리는 [앱을 위한 키퍼들](#앱을 위한 키퍼들) 절에서 우리 모듈이 두 개의 다른 모듈인 `bank` 모듈과 `stake` 모듈에 대한 액세스를 필요로한다는 것을 분명히했습니다. 또한 기본적인 계정 기능을 위한 `auth` 모듈이 필요합니다. 마지막으로 우리가 사용하는 각 모듈의 저장소를 선언하기 위해 메인 멀티 스토어에 액세스해야합니다.

#### 앱 명령어들

새로 생성된 명령을 어플리케이션에 추가해야합니다. 이렇게하려면 루트 디렉토리의 `cmd` 폴더로 가십시오:

```bash
// 디렉토리의 루트 레벨에서
cd cmd
```
`simplegovd`는 서버 데몬을 실행하기위한 명령을 저장하는 폴더이며, `simplegovcli`는 어플리케이션의 명령을 정의합니다.

##### CLI

어플리케이션과 상호 작용하려면 `simple_governance` 모듈의 명령을`simpleGov` 어플리케이션과 사전 빌드된 SDK 명령에 추가하십시오.

```go
//  cmd/simplegovcli/main.go
...
	rootCmd.AddCommand(
		client.GetCommands(
			simplegovcmd.GetCmdQueryProposal("proposals", cdc),
			simplegovcmd.GetCmdQueryProposals("proposals", cdc),
			simplegovcmd.GetCmdQueryProposalVotes("proposals", cdc),
			simplegovcmd.GetCmdQueryProposalVote("proposals", cdc),
)...)
	rootCmd.AddCommand(
		client.PostCommands(
			simplegovcmd.PostCmdPropose(cdc),
			simplegovcmd.PostCmdVote(cdc),
)...)
...
```

##### 데몬 서버

`simplegovd` 명령은 데몬 서버를 백그라운드 프로세스로 실행합니다. 먼저 `utils` 함수를 만들어 봅시다:

```go
//  cmd/simplegovd/main.go
// SimpleGovAppInit 초기 매개 변수
var SimpleGovAppInit = server.AppInit{
	AppGenState: SimpleGovAppGenState,
	AppGenTx:    server.SimpleAppGenTx,
}

// SimpleGovAppGenState는 app_state를 설정하고 simpleGov 앱 상태를 추가합니다.
func SimpleGovAppGenState(cdc *wire.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {
	appState, err = server.SimpleAppGenState(cdc, appGenTxs)
	if err != nil {
		return
}
	return
}

func newApp(logger log.Logger, db dbm.DB) abci.Application {
	return app.NewSimpleGovApp(logger, db)
}

func exportAppState(logger log.Logger, db dbm.DB) (json.RawMessage, error) {
	dapp := app.NewSimpleGovApp(logger, db)
	return dapp.ExportAppStateJSON()
}
```

이제 `main()` 함수 내에서 데몬서버에 대한 명령을 정의해 보겠습니다.

```go
//  cmd/simplegovd/main.go
func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "simplegovd",
		Short:             "Simple Governance Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
}

	server.AddCommands(ctx, cdc, rootCmd, SimpleGovAppInit,
		server.ConstructAppCreator(newApp, "simplegov"),
		server.ConstructAppExporter(exportAppState, "simplegov"))

	// 플래그를 준비하고 추가합니다.
	rootDir := os.ExpandEnv("$HOME/.simplegovd")
	executor := cli.PrepareBaseCmd(rootCmd, "BC", rootDir)
	executor.Execute()
}
```

##### Makefile

[Makefile](https://en.wikipedia.org/wiki/Makefile)은 대상과 레시피가 있는 일련의 규칙을 정의하여 Go 프로그램을 컴파일합니다. 어플리케이션 명령을 추가해야합니다.

```
##### Makefile
build_examples:
ifeq ($(OS),Windows_NT)
...
	go build $(BUILD_FLAGS) -o build/simplegovd.exe ./examples/simpleGov/cmd/simplegovd
	go build $(BUILD_FLAGS) -o build/simplegovcli.exe ./examples/simpleGov/cmd/simplegovcli
else
...
	go build $(BUILD_FLAGS) -o build/simplegovd ./examples/simpleGov/cmd/simplegovd
	go build $(BUILD_FLAGS) -o build/simplegovcli ./examples/simpleGov/cmd/simplegovcli
endif
...
install_examples:
    ...
	go install $(BUILD_FLAGS) ./examples/simpleGov/cmd/simplegovd
	go install $(BUILD_FLAGS) ./examples/simpleGov/cmd/simplegovcli
```

#### App 생성자

이제 어플리케이션용 생성자를 정의해야합니다.

```go
func NewSimpleGovApp(logger log.Logger, db dbm.DB) *SimpleGovApp
```

이 기능에서 우리는:

- 코덱 생성

```go
var cdc = MakeCodec()
```

- 우리의 어플리케이션을 인스턴스화하십시오. 여기에는 각 하위 노드에 액세스 하기 위한 키 생성이 포함됩니다.

```go
// 어플리케이션 객체를 만듭니다.
    var app = &SimpleGovApp{
        BaseApp:              bam.NewBaseApp(appName, cdc, logger, db),
        cdc:                  cdc,
        capKeyMainStore:      sdk.NewKVStoreKey("main"),
        capKeyAccountStore:   sdk.NewKVStoreKey("acc"),
        capKeyStakingStore:   sdk.NewKVStoreKey("stake"),
        capKeySimpleGovStore: sdk.NewKVStoreKey("simpleGov"),
    }
```

- 키퍼를 인스턴스화합니다. 키퍼는 일반적으로 다른 모듈의 키퍼에 액세스해야합니다. 이 경우, 필요한 기능에 대한 키퍼 인스턴스만 전달해야합니다. 키퍼가 다른 모듈의 저장소에서 읽기만하면 읽기 전용 키퍼를 전달해야합니다.

```go
app.coinKeeper = bank.NewKeeper(app.accountMapper)
app.stakeKeeper = simplestake.NewKeeper(app.capKeyStakingStore, app.coinKeeper,app.RegisterCodespace(simplestake.DefaultCodespace))
app.simpleGovKeeper = simpleGov.NewKeeper(app.capKeySimpleGovStore, app.coinKeeper, app.stakeKeeper, app.RegisterCodespace(simpleGov.DefaultCodespace))
```

- 핸들러를 선언하십시오.

```go
app.Router().
        AddRoute("bank", bank.NewHandler(app.coinKeeper)).
        AddRoute("simplestake", simplestake.NewHandler(app.stakeKeeper)).
        AddRoute("simpleGov", simpleGov.NewHandler(app.simpleGovKeeper))
```

- 어플리케이션을 초기화하십시오.

```go
// BaseApp 초기화.
    app.MountStoresIAVL(app.capKeyMainStore, app.capKeyAccountStore, app.capKeySimpleGovStore, app.capKeyStakingStore)
    app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))
    err := app.LoadLatestVersion(app.capKeyMainStore)
    if err != nil {
        cmn.Exit(err.Error())
    }
    return app
```

#### 앱 코덱

마지막으로, 우리는 `MakeCodec()` 함수를 정의하고 다양한 모듈의 구체적인 타입과 인터페이스를 등록해야합니다.

```go
func MakeCodec() *wire.Codec {
    var cdc = wire.NewCodec()
    wire.RegisterCrypto(cdc) // 암호를 등록.
    sdk.RegisterWire(cdc)    // Msgs 등록
    bank.RegisterWire(cdc)
    simplestake.RegisterWire(cdc)
    simpleGov.RegisterWire(cdc)

    // AppAccount 등록
    cdc.RegisterInterface((*auth.Account)(nil), nil)
    cdc.RegisterConcrete(&types.AppAccount{}, "simpleGov/Account", nil)
    return cdc
}
```

### 앱 실행하기

#### 설치

일단 어플리케이션을 마무리 했다면 `go get`을 사용하여 어플리케이션을 설치하십시오. 다음의 명령은 `simpleGov` 어플리케이션뿐만 아니라 미리 만들어진 모듈과 SDK 예제를 설치합니다:

```bash
go get github.com/<your_username>/cosmos-sdk
cd $GOPATH/src/github.com/<your_username>/cosmos-sdk
make get_vendor_deps
make install
make install_examples
```

다음을 입력하여 어플리케이션이 올바르게 설치되었는지 확인하십시오:

```bash
simplegovcli -h
simplegovd -h
```

#### 제안서 제출

CLI를 사용하여 새로운 제안서 작성:

```bash
simplegovcli propose --title="Voting Period update" --description="Should we change the proposal voting period to 3 weeks?" --deposit=300Atoms
```

새로 생성된 제안서의 세부 정보 얻기:

```bash
simplegovcli proposal 1
```

기존의 모든 제안서를 확인할 수도 있습니다:

```bash
simplegovcli proposals --active=true
```

#### 기존 제안서에 투표하기

생성된 제안서에 투표해 봅시다:

```bash
simplegovcli vote --proposal-id=1 --option="No"
```

투표한 표에서 옵션 값을 얻기:

```bash
simplegovcli proposal-vote 1 <your_address>
```

또한 제안서에 투표한 모든 표들을 확인할 수 있습니다:

```bash
simplegovcli proposals-votes 1
```

### Testnet

WIP

### 결론

축하합니다! Cosmos-SDK로 첫 번째 어플리케이션과 모듈을 성공적으로 만들었습니다. 이 자습서 또는 SDK 개발에 관한 질문이 있으시면 공식 커뮤니케이션 채널을 통해 문의하십시오:

- [Cosmos-SDK Riot Channel](https://riot.im/app/#/room/#cosmos-sdk:matrix.org)
- [텔레그램](https://t.me/cosmosproject)
