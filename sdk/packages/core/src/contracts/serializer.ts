export interface TransportSerializer<TInput = unknown, TOutput = string> {
  contentType: string;
  serialize(input: TInput): TOutput;
  deserialize(output: TOutput): TInput;
}
