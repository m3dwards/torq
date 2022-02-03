/* eslint-disable */
import Long from "long";
import { grpc } from "@improbable-eng/grpc-web";
import _m0 from "protobufjs/minimal";
import { BrowserHeaders } from "browser-headers";

export const protobufPackage = "torqrpc";

export interface ChannelFlowRequest {
  /** From what date/time (unix timestamp) */
  fromTime: Long;
  /** To what date/time (unix timestamp) */
  toTime: Long;
  /** TODO: Add "repeated" here to request multiple channels at once. */
  chanIds: Long[];
}

export interface ChannelFlow {
  /** What channel ID's the flow is for */
  chanIds: Long[];
  /** From what date/time (unix timestamp) */
  fromTime: Long;
  /** To what date/time (unix timestamp) */
  toTime: Long;
  /** Fees earned by other channels using this channels inbound liquidity. */
  feeIn: Long;
  /** Fees earned by this channels outbound liquidity */
  feeOut: Long;
  /** Amount inbound */
  amtIn: Long;
  /** Amount outbound */
  amtOut: Long;
  /** Number of forwards inbound */
  countIn: Long;
  /** Number of forwards outbound */
  countOut: Long;
}

export interface AggregatedForwards {
  fromTime: Long;
  toTime: Long;
  /** The incoming channel ID that carried the HTLC that created the circuit. */
  chanIdIn: Long;
  /**
   * The outgoing channel ID that carried the preimage that completed the
   * circuit.
   */
  chanIdOut: Long;
  /** The total fee (in satoshis) that this payment circuit carried. */
  fee: Long;
  /**
   * The total amount (in satoshis) of the incoming HTLC that created half
   * the circuit.
   */
  amtIn: Long;
  /**
   * The total amount (in satoshis) of the outgoing HTLC that created the
   * second half of the circuit.
   */
  amtOut: Long;
  countIn: Long;
  countOut: Long;
}

export interface Forwards {
  forwards: Forward[];
}

export interface Forward {
  /** The microseconds' version of TimestampNs, used by TimescaleDB */
  time: Long;
  /**
   * The number of nanoseconds elapsed since January 1, 1970 UTC when this
   * circuit was completed.
   */
  timeNs: Long;
  /** The incoming channel ID that carried the HTLC that created the circuit. */
  chanIdIn: Long;
  /**
   * The outgoing channel ID that carried the preimage that completed the
   * circuit.
   */
  chanIdOut: Long;
  /** The total fee (in satoshis) that this payment circuit carried. */
  fee: Long;
  /**
   * The total amount (in satoshis) of the incoming HTLC that created half
   * the circuit.
   */
  amtIn: Long;
  /**
   * The total amount (in satoshis) of the outgoing HTLC that created the
   * second half of the circuit.
   */
  amtOut: Long;
}

export interface ForwardsRequest {}

const baseChannelFlowRequest: object = {
  fromTime: Long.ZERO,
  toTime: Long.ZERO,
  chanIds: Long.UZERO,
};

export const ChannelFlowRequest = {
  encode(
    message: ChannelFlowRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.fromTime.isZero()) {
      writer.uint32(8).int64(message.fromTime);
    }
    if (!message.toTime.isZero()) {
      writer.uint32(16).int64(message.toTime);
    }
    writer.uint32(26).fork();
    for (const v of message.chanIds) {
      writer.uint64(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChannelFlowRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseChannelFlowRequest } as ChannelFlowRequest;
    message.chanIds = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fromTime = reader.int64() as Long;
          break;
        case 2:
          message.toTime = reader.int64() as Long;
          break;
        case 3:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.chanIds.push(reader.uint64() as Long);
            }
          } else {
            message.chanIds.push(reader.uint64() as Long);
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ChannelFlowRequest {
    const message = { ...baseChannelFlowRequest } as ChannelFlowRequest;
    message.fromTime =
      object.fromTime !== undefined && object.fromTime !== null
        ? Long.fromString(object.fromTime)
        : Long.ZERO;
    message.toTime =
      object.toTime !== undefined && object.toTime !== null
        ? Long.fromString(object.toTime)
        : Long.ZERO;
    message.chanIds = (object.chanIds ?? []).map((e: any) =>
      Long.fromString(e)
    );
    return message;
  },

  toJSON(message: ChannelFlowRequest): unknown {
    const obj: any = {};
    message.fromTime !== undefined &&
      (obj.fromTime = (message.fromTime || Long.ZERO).toString());
    message.toTime !== undefined &&
      (obj.toTime = (message.toTime || Long.ZERO).toString());
    if (message.chanIds) {
      obj.chanIds = message.chanIds.map((e) => (e || Long.UZERO).toString());
    } else {
      obj.chanIds = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ChannelFlowRequest>, I>>(
    object: I
  ): ChannelFlowRequest {
    const message = { ...baseChannelFlowRequest } as ChannelFlowRequest;
    message.fromTime =
      object.fromTime !== undefined && object.fromTime !== null
        ? Long.fromValue(object.fromTime)
        : Long.ZERO;
    message.toTime =
      object.toTime !== undefined && object.toTime !== null
        ? Long.fromValue(object.toTime)
        : Long.ZERO;
    message.chanIds = object.chanIds?.map((e) => Long.fromValue(e)) || [];
    return message;
  },
};

const baseChannelFlow: object = {
  chanIds: Long.UZERO,
  fromTime: Long.ZERO,
  toTime: Long.ZERO,
  feeIn: Long.UZERO,
  feeOut: Long.UZERO,
  amtIn: Long.UZERO,
  amtOut: Long.UZERO,
  countIn: Long.ZERO,
  countOut: Long.ZERO,
};

export const ChannelFlow = {
  encode(
    message: ChannelFlow,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    writer.uint32(10).fork();
    for (const v of message.chanIds) {
      writer.uint64(v);
    }
    writer.ldelim();
    if (!message.fromTime.isZero()) {
      writer.uint32(16).int64(message.fromTime);
    }
    if (!message.toTime.isZero()) {
      writer.uint32(24).int64(message.toTime);
    }
    if (!message.feeIn.isZero()) {
      writer.uint32(32).uint64(message.feeIn);
    }
    if (!message.feeOut.isZero()) {
      writer.uint32(40).uint64(message.feeOut);
    }
    if (!message.amtIn.isZero()) {
      writer.uint32(48).uint64(message.amtIn);
    }
    if (!message.amtOut.isZero()) {
      writer.uint32(56).uint64(message.amtOut);
    }
    if (!message.countIn.isZero()) {
      writer.uint32(64).int64(message.countIn);
    }
    if (!message.countOut.isZero()) {
      writer.uint32(72).int64(message.countOut);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChannelFlow {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseChannelFlow } as ChannelFlow;
    message.chanIds = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.chanIds.push(reader.uint64() as Long);
            }
          } else {
            message.chanIds.push(reader.uint64() as Long);
          }
          break;
        case 2:
          message.fromTime = reader.int64() as Long;
          break;
        case 3:
          message.toTime = reader.int64() as Long;
          break;
        case 4:
          message.feeIn = reader.uint64() as Long;
          break;
        case 5:
          message.feeOut = reader.uint64() as Long;
          break;
        case 6:
          message.amtIn = reader.uint64() as Long;
          break;
        case 7:
          message.amtOut = reader.uint64() as Long;
          break;
        case 8:
          message.countIn = reader.int64() as Long;
          break;
        case 9:
          message.countOut = reader.int64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ChannelFlow {
    const message = { ...baseChannelFlow } as ChannelFlow;
    message.chanIds = (object.chanIds ?? []).map((e: any) =>
      Long.fromString(e)
    );
    message.fromTime =
      object.fromTime !== undefined && object.fromTime !== null
        ? Long.fromString(object.fromTime)
        : Long.ZERO;
    message.toTime =
      object.toTime !== undefined && object.toTime !== null
        ? Long.fromString(object.toTime)
        : Long.ZERO;
    message.feeIn =
      object.feeIn !== undefined && object.feeIn !== null
        ? Long.fromString(object.feeIn)
        : Long.UZERO;
    message.feeOut =
      object.feeOut !== undefined && object.feeOut !== null
        ? Long.fromString(object.feeOut)
        : Long.UZERO;
    message.amtIn =
      object.amtIn !== undefined && object.amtIn !== null
        ? Long.fromString(object.amtIn)
        : Long.UZERO;
    message.amtOut =
      object.amtOut !== undefined && object.amtOut !== null
        ? Long.fromString(object.amtOut)
        : Long.UZERO;
    message.countIn =
      object.countIn !== undefined && object.countIn !== null
        ? Long.fromString(object.countIn)
        : Long.ZERO;
    message.countOut =
      object.countOut !== undefined && object.countOut !== null
        ? Long.fromString(object.countOut)
        : Long.ZERO;
    return message;
  },

  toJSON(message: ChannelFlow): unknown {
    const obj: any = {};
    if (message.chanIds) {
      obj.chanIds = message.chanIds.map((e) => (e || Long.UZERO).toString());
    } else {
      obj.chanIds = [];
    }
    message.fromTime !== undefined &&
      (obj.fromTime = (message.fromTime || Long.ZERO).toString());
    message.toTime !== undefined &&
      (obj.toTime = (message.toTime || Long.ZERO).toString());
    message.feeIn !== undefined &&
      (obj.feeIn = (message.feeIn || Long.UZERO).toString());
    message.feeOut !== undefined &&
      (obj.feeOut = (message.feeOut || Long.UZERO).toString());
    message.amtIn !== undefined &&
      (obj.amtIn = (message.amtIn || Long.UZERO).toString());
    message.amtOut !== undefined &&
      (obj.amtOut = (message.amtOut || Long.UZERO).toString());
    message.countIn !== undefined &&
      (obj.countIn = (message.countIn || Long.ZERO).toString());
    message.countOut !== undefined &&
      (obj.countOut = (message.countOut || Long.ZERO).toString());
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ChannelFlow>, I>>(
    object: I
  ): ChannelFlow {
    const message = { ...baseChannelFlow } as ChannelFlow;
    message.chanIds = object.chanIds?.map((e) => Long.fromValue(e)) || [];
    message.fromTime =
      object.fromTime !== undefined && object.fromTime !== null
        ? Long.fromValue(object.fromTime)
        : Long.ZERO;
    message.toTime =
      object.toTime !== undefined && object.toTime !== null
        ? Long.fromValue(object.toTime)
        : Long.ZERO;
    message.feeIn =
      object.feeIn !== undefined && object.feeIn !== null
        ? Long.fromValue(object.feeIn)
        : Long.UZERO;
    message.feeOut =
      object.feeOut !== undefined && object.feeOut !== null
        ? Long.fromValue(object.feeOut)
        : Long.UZERO;
    message.amtIn =
      object.amtIn !== undefined && object.amtIn !== null
        ? Long.fromValue(object.amtIn)
        : Long.UZERO;
    message.amtOut =
      object.amtOut !== undefined && object.amtOut !== null
        ? Long.fromValue(object.amtOut)
        : Long.UZERO;
    message.countIn =
      object.countIn !== undefined && object.countIn !== null
        ? Long.fromValue(object.countIn)
        : Long.ZERO;
    message.countOut =
      object.countOut !== undefined && object.countOut !== null
        ? Long.fromValue(object.countOut)
        : Long.ZERO;
    return message;
  },
};

const baseAggregatedForwards: object = {
  fromTime: Long.ZERO,
  toTime: Long.ZERO,
  chanIdIn: Long.UZERO,
  chanIdOut: Long.UZERO,
  fee: Long.UZERO,
  amtIn: Long.UZERO,
  amtOut: Long.UZERO,
  countIn: Long.ZERO,
  countOut: Long.ZERO,
};

export const AggregatedForwards = {
  encode(
    message: AggregatedForwards,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.fromTime.isZero()) {
      writer.uint32(8).int64(message.fromTime);
    }
    if (!message.toTime.isZero()) {
      writer.uint32(16).int64(message.toTime);
    }
    if (!message.chanIdIn.isZero()) {
      writer.uint32(24).uint64(message.chanIdIn);
    }
    if (!message.chanIdOut.isZero()) {
      writer.uint32(32).uint64(message.chanIdOut);
    }
    if (!message.fee.isZero()) {
      writer.uint32(40).uint64(message.fee);
    }
    if (!message.amtIn.isZero()) {
      writer.uint32(48).uint64(message.amtIn);
    }
    if (!message.amtOut.isZero()) {
      writer.uint32(56).uint64(message.amtOut);
    }
    if (!message.countIn.isZero()) {
      writer.uint32(64).int64(message.countIn);
    }
    if (!message.countOut.isZero()) {
      writer.uint32(72).int64(message.countOut);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AggregatedForwards {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAggregatedForwards } as AggregatedForwards;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fromTime = reader.int64() as Long;
          break;
        case 2:
          message.toTime = reader.int64() as Long;
          break;
        case 3:
          message.chanIdIn = reader.uint64() as Long;
          break;
        case 4:
          message.chanIdOut = reader.uint64() as Long;
          break;
        case 5:
          message.fee = reader.uint64() as Long;
          break;
        case 6:
          message.amtIn = reader.uint64() as Long;
          break;
        case 7:
          message.amtOut = reader.uint64() as Long;
          break;
        case 8:
          message.countIn = reader.int64() as Long;
          break;
        case 9:
          message.countOut = reader.int64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AggregatedForwards {
    const message = { ...baseAggregatedForwards } as AggregatedForwards;
    message.fromTime =
      object.fromTime !== undefined && object.fromTime !== null
        ? Long.fromString(object.fromTime)
        : Long.ZERO;
    message.toTime =
      object.toTime !== undefined && object.toTime !== null
        ? Long.fromString(object.toTime)
        : Long.ZERO;
    message.chanIdIn =
      object.chanIdIn !== undefined && object.chanIdIn !== null
        ? Long.fromString(object.chanIdIn)
        : Long.UZERO;
    message.chanIdOut =
      object.chanIdOut !== undefined && object.chanIdOut !== null
        ? Long.fromString(object.chanIdOut)
        : Long.UZERO;
    message.fee =
      object.fee !== undefined && object.fee !== null
        ? Long.fromString(object.fee)
        : Long.UZERO;
    message.amtIn =
      object.amtIn !== undefined && object.amtIn !== null
        ? Long.fromString(object.amtIn)
        : Long.UZERO;
    message.amtOut =
      object.amtOut !== undefined && object.amtOut !== null
        ? Long.fromString(object.amtOut)
        : Long.UZERO;
    message.countIn =
      object.countIn !== undefined && object.countIn !== null
        ? Long.fromString(object.countIn)
        : Long.ZERO;
    message.countOut =
      object.countOut !== undefined && object.countOut !== null
        ? Long.fromString(object.countOut)
        : Long.ZERO;
    return message;
  },

  toJSON(message: AggregatedForwards): unknown {
    const obj: any = {};
    message.fromTime !== undefined &&
      (obj.fromTime = (message.fromTime || Long.ZERO).toString());
    message.toTime !== undefined &&
      (obj.toTime = (message.toTime || Long.ZERO).toString());
    message.chanIdIn !== undefined &&
      (obj.chanIdIn = (message.chanIdIn || Long.UZERO).toString());
    message.chanIdOut !== undefined &&
      (obj.chanIdOut = (message.chanIdOut || Long.UZERO).toString());
    message.fee !== undefined &&
      (obj.fee = (message.fee || Long.UZERO).toString());
    message.amtIn !== undefined &&
      (obj.amtIn = (message.amtIn || Long.UZERO).toString());
    message.amtOut !== undefined &&
      (obj.amtOut = (message.amtOut || Long.UZERO).toString());
    message.countIn !== undefined &&
      (obj.countIn = (message.countIn || Long.ZERO).toString());
    message.countOut !== undefined &&
      (obj.countOut = (message.countOut || Long.ZERO).toString());
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AggregatedForwards>, I>>(
    object: I
  ): AggregatedForwards {
    const message = { ...baseAggregatedForwards } as AggregatedForwards;
    message.fromTime =
      object.fromTime !== undefined && object.fromTime !== null
        ? Long.fromValue(object.fromTime)
        : Long.ZERO;
    message.toTime =
      object.toTime !== undefined && object.toTime !== null
        ? Long.fromValue(object.toTime)
        : Long.ZERO;
    message.chanIdIn =
      object.chanIdIn !== undefined && object.chanIdIn !== null
        ? Long.fromValue(object.chanIdIn)
        : Long.UZERO;
    message.chanIdOut =
      object.chanIdOut !== undefined && object.chanIdOut !== null
        ? Long.fromValue(object.chanIdOut)
        : Long.UZERO;
    message.fee =
      object.fee !== undefined && object.fee !== null
        ? Long.fromValue(object.fee)
        : Long.UZERO;
    message.amtIn =
      object.amtIn !== undefined && object.amtIn !== null
        ? Long.fromValue(object.amtIn)
        : Long.UZERO;
    message.amtOut =
      object.amtOut !== undefined && object.amtOut !== null
        ? Long.fromValue(object.amtOut)
        : Long.UZERO;
    message.countIn =
      object.countIn !== undefined && object.countIn !== null
        ? Long.fromValue(object.countIn)
        : Long.ZERO;
    message.countOut =
      object.countOut !== undefined && object.countOut !== null
        ? Long.fromValue(object.countOut)
        : Long.ZERO;
    return message;
  },
};

const baseForwards: object = {};

export const Forwards = {
  encode(
    message: Forwards,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.forwards) {
      Forward.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Forwards {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseForwards } as Forwards;
    message.forwards = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.forwards.push(Forward.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Forwards {
    const message = { ...baseForwards } as Forwards;
    message.forwards = (object.forwards ?? []).map((e: any) =>
      Forward.fromJSON(e)
    );
    return message;
  },

  toJSON(message: Forwards): unknown {
    const obj: any = {};
    if (message.forwards) {
      obj.forwards = message.forwards.map((e) =>
        e ? Forward.toJSON(e) : undefined
      );
    } else {
      obj.forwards = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Forwards>, I>>(object: I): Forwards {
    const message = { ...baseForwards } as Forwards;
    message.forwards =
      object.forwards?.map((e) => Forward.fromPartial(e)) || [];
    return message;
  },
};

const baseForward: object = {
  time: Long.ZERO,
  timeNs: Long.ZERO,
  chanIdIn: Long.UZERO,
  chanIdOut: Long.UZERO,
  fee: Long.UZERO,
  amtIn: Long.UZERO,
  amtOut: Long.UZERO,
};

export const Forward = {
  encode(
    message: Forward,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (!message.time.isZero()) {
      writer.uint32(8).int64(message.time);
    }
    if (!message.timeNs.isZero()) {
      writer.uint32(16).int64(message.timeNs);
    }
    if (!message.chanIdIn.isZero()) {
      writer.uint32(24).uint64(message.chanIdIn);
    }
    if (!message.chanIdOut.isZero()) {
      writer.uint32(32).uint64(message.chanIdOut);
    }
    if (!message.fee.isZero()) {
      writer.uint32(40).uint64(message.fee);
    }
    if (!message.amtIn.isZero()) {
      writer.uint32(48).uint64(message.amtIn);
    }
    if (!message.amtOut.isZero()) {
      writer.uint32(56).uint64(message.amtOut);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Forward {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseForward } as Forward;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.time = reader.int64() as Long;
          break;
        case 2:
          message.timeNs = reader.int64() as Long;
          break;
        case 3:
          message.chanIdIn = reader.uint64() as Long;
          break;
        case 4:
          message.chanIdOut = reader.uint64() as Long;
          break;
        case 5:
          message.fee = reader.uint64() as Long;
          break;
        case 6:
          message.amtIn = reader.uint64() as Long;
          break;
        case 7:
          message.amtOut = reader.uint64() as Long;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Forward {
    const message = { ...baseForward } as Forward;
    message.time =
      object.time !== undefined && object.time !== null
        ? Long.fromString(object.time)
        : Long.ZERO;
    message.timeNs =
      object.timeNs !== undefined && object.timeNs !== null
        ? Long.fromString(object.timeNs)
        : Long.ZERO;
    message.chanIdIn =
      object.chanIdIn !== undefined && object.chanIdIn !== null
        ? Long.fromString(object.chanIdIn)
        : Long.UZERO;
    message.chanIdOut =
      object.chanIdOut !== undefined && object.chanIdOut !== null
        ? Long.fromString(object.chanIdOut)
        : Long.UZERO;
    message.fee =
      object.fee !== undefined && object.fee !== null
        ? Long.fromString(object.fee)
        : Long.UZERO;
    message.amtIn =
      object.amtIn !== undefined && object.amtIn !== null
        ? Long.fromString(object.amtIn)
        : Long.UZERO;
    message.amtOut =
      object.amtOut !== undefined && object.amtOut !== null
        ? Long.fromString(object.amtOut)
        : Long.UZERO;
    return message;
  },

  toJSON(message: Forward): unknown {
    const obj: any = {};
    message.time !== undefined &&
      (obj.time = (message.time || Long.ZERO).toString());
    message.timeNs !== undefined &&
      (obj.timeNs = (message.timeNs || Long.ZERO).toString());
    message.chanIdIn !== undefined &&
      (obj.chanIdIn = (message.chanIdIn || Long.UZERO).toString());
    message.chanIdOut !== undefined &&
      (obj.chanIdOut = (message.chanIdOut || Long.UZERO).toString());
    message.fee !== undefined &&
      (obj.fee = (message.fee || Long.UZERO).toString());
    message.amtIn !== undefined &&
      (obj.amtIn = (message.amtIn || Long.UZERO).toString());
    message.amtOut !== undefined &&
      (obj.amtOut = (message.amtOut || Long.UZERO).toString());
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Forward>, I>>(object: I): Forward {
    const message = { ...baseForward } as Forward;
    message.time =
      object.time !== undefined && object.time !== null
        ? Long.fromValue(object.time)
        : Long.ZERO;
    message.timeNs =
      object.timeNs !== undefined && object.timeNs !== null
        ? Long.fromValue(object.timeNs)
        : Long.ZERO;
    message.chanIdIn =
      object.chanIdIn !== undefined && object.chanIdIn !== null
        ? Long.fromValue(object.chanIdIn)
        : Long.UZERO;
    message.chanIdOut =
      object.chanIdOut !== undefined && object.chanIdOut !== null
        ? Long.fromValue(object.chanIdOut)
        : Long.UZERO;
    message.fee =
      object.fee !== undefined && object.fee !== null
        ? Long.fromValue(object.fee)
        : Long.UZERO;
    message.amtIn =
      object.amtIn !== undefined && object.amtIn !== null
        ? Long.fromValue(object.amtIn)
        : Long.UZERO;
    message.amtOut =
      object.amtOut !== undefined && object.amtOut !== null
        ? Long.fromValue(object.amtOut)
        : Long.UZERO;
    return message;
  },
};

const baseForwardsRequest: object = {};

export const ForwardsRequest = {
  encode(
    _: ForwardsRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ForwardsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseForwardsRequest } as ForwardsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ForwardsRequest {
    const message = { ...baseForwardsRequest } as ForwardsRequest;
    return message;
  },

  toJSON(_: ForwardsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ForwardsRequest>, I>>(
    _: I
  ): ForwardsRequest {
    const message = { ...baseForwardsRequest } as ForwardsRequest;
    return message;
  },
};

export interface torqrpc {
  GetForwards(
    request: DeepPartial<ForwardsRequest>,
    metadata?: grpc.Metadata
  ): Promise<Forwards>;
  GetChannelFlow(
    request: DeepPartial<ChannelFlowRequest>,
    metadata?: grpc.Metadata
  ): Promise<ChannelFlow>;
}

export class torqrpcClientImpl implements torqrpc {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.GetForwards = this.GetForwards.bind(this);
    this.GetChannelFlow = this.GetChannelFlow.bind(this);
  }

  GetForwards(
    request: DeepPartial<ForwardsRequest>,
    metadata?: grpc.Metadata
  ): Promise<Forwards> {
    return this.rpc.unary(
      torqrpcGetForwardsDesc,
      ForwardsRequest.fromPartial(request),
      metadata
    );
  }

  GetChannelFlow(
    request: DeepPartial<ChannelFlowRequest>,
    metadata?: grpc.Metadata
  ): Promise<ChannelFlow> {
    return this.rpc.unary(
      torqrpcGetChannelFlowDesc,
      ChannelFlowRequest.fromPartial(request),
      metadata
    );
  }
}

export const torqrpcDesc = {
  serviceName: "torqrpc.torqrpc",
};

export const torqrpcGetForwardsDesc: UnaryMethodDefinitionish = {
  methodName: "GetForwards",
  service: torqrpcDesc,
  requestStream: false,
  responseStream: false,
  requestType: {
    serializeBinary() {
      return ForwardsRequest.encode(this).finish();
    },
  } as any,
  responseType: {
    deserializeBinary(data: Uint8Array) {
      return {
        ...Forwards.decode(data),
        toObject() {
          return this;
        },
      };
    },
  } as any,
};

export const torqrpcGetChannelFlowDesc: UnaryMethodDefinitionish = {
  methodName: "GetChannelFlow",
  service: torqrpcDesc,
  requestStream: false,
  responseStream: false,
  requestType: {
    serializeBinary() {
      return ChannelFlowRequest.encode(this).finish();
    },
  } as any,
  responseType: {
    deserializeBinary(data: Uint8Array) {
      return {
        ...ChannelFlow.decode(data),
        toObject() {
          return this;
        },
      };
    },
  } as any,
};

interface UnaryMethodDefinitionishR
  extends grpc.UnaryMethodDefinition<any, any> {
  requestStream: any;
  responseStream: any;
}

type UnaryMethodDefinitionish = UnaryMethodDefinitionishR;

interface Rpc {
  unary<T extends UnaryMethodDefinitionish>(
    methodDesc: T,
    request: any,
    metadata: grpc.Metadata | undefined
  ): Promise<any>;
}

export class GrpcWebImpl {
  private host: string;
  private options: {
    transport?: grpc.TransportFactory;

    debug?: boolean;
    metadata?: grpc.Metadata;
  };

  constructor(
    host: string,
    options: {
      transport?: grpc.TransportFactory;

      debug?: boolean;
      metadata?: grpc.Metadata;
    }
  ) {
    this.host = host;
    this.options = options;
  }

  unary<T extends UnaryMethodDefinitionish>(
    methodDesc: T,
    _request: any,
    metadata: grpc.Metadata | undefined
  ): Promise<any> {
    const request = { ..._request, ...methodDesc.requestType };
    const maybeCombinedMetadata =
      metadata && this.options.metadata
        ? new BrowserHeaders({
            ...this.options?.metadata.headersMap,
            ...metadata?.headersMap,
          })
        : metadata || this.options.metadata;
    return new Promise((resolve, reject) => {
      grpc.unary(methodDesc, {
        request,
        host: this.host,
        metadata: maybeCombinedMetadata,
        transport: this.options.transport,
        debug: this.options.debug,
        onEnd: function (response) {
          if (response.status === grpc.Code.OK) {
            resolve(response.message);
          } else {
            const err = new Error(response.statusMessage) as any;
            err.code = response.status;
            err.metadata = response.trailers;
            reject(err);
          }
        },
      });
    });
  }
}

type Builtin =
  | Date
  | Function
  | Uint8Array
  | string
  | number
  | boolean
  | undefined;

export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Long
  ? string | number | Long
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin
  ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & Record<
        Exclude<keyof I, KeysOfUnion<P>>,
        never
      >;

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}
