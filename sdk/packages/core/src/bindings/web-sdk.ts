import { ServifySDK } from '../sdk';
import type { ServifyConfig } from '../types';

export type WebServifyConfig = ServifyConfig;
export type WebServifyClient = ServifySDK;

export function createWebServifySDK(config: WebServifyConfig): WebServifyClient {
  return new ServifySDK(config);
}
