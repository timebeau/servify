export interface PushTokenRegistration {
  token: string;
  platform: 'ios' | 'android';
  deviceId: string;
  environment?: 'sandbox' | 'production';
}

export interface PushTokenRegistrar {
  register(input: PushTokenRegistration): Promise<void>;
  unregister(deviceId: string): Promise<void>;
}
