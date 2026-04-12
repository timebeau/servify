import { useCallback, useEffect, useState } from 'react';
import type { RemoteAssistStartOptions, RemoteAssistState } from '@servify/core';
import { useServify } from './ServifyProvider';

export interface UseRemoteAssistReturn {
  state: RemoteAssistState;
  isActive: boolean;
  error: Error | null;
  remoteStream: MediaStream | null;
  startRemoteAssist: (options?: RemoteAssistStartOptions) => Promise<void>;
  acceptRemoteAnswer: (answer: RTCSessionDescriptionInit) => Promise<void>;
  addRemoteIce: (candidate: RTCIceCandidateInit) => Promise<void>;
  endRemoteAssist: () => Promise<void>;
}

export function useRemoteAssist(): UseRemoteAssistReturn {
  const { sdk } = useServify();
  const [state, setState] = useState<RemoteAssistState>('idle');
  const [error, setError] = useState<Error | null>(null);
  const [remoteStream, setRemoteStream] = useState<MediaStream | null>(null);

  const startRemoteAssist = useCallback(async (options?: RemoteAssistStartOptions) => {
    if (!sdk) {
      throw new Error('SDK not initialized');
    }

    try {
      setError(null);
      await sdk.startRemoteAssist(options);
    } catch (err) {
      setError(err as Error);
      throw err;
    }
  }, [sdk]);

  const acceptRemoteAnswer = useCallback(async (answer: RTCSessionDescriptionInit) => {
    if (!sdk) {
      throw new Error('SDK not initialized');
    }

    try {
      setError(null);
      await sdk.acceptRemoteAnswer(answer);
    } catch (err) {
      setError(err as Error);
      throw err;
    }
  }, [sdk]);

  const addRemoteIce = useCallback(async (candidate: RTCIceCandidateInit) => {
    if (!sdk) {
      throw new Error('SDK not initialized');
    }

    try {
      setError(null);
      await sdk.addRemoteIce(candidate);
    } catch (err) {
      setError(err as Error);
      throw err;
    }
  }, [sdk]);

  const endRemoteAssist = useCallback(async () => {
    if (!sdk) {
      throw new Error('SDK not initialized');
    }

    try {
      setError(null);
      await sdk.endRemoteAssist();
      setRemoteStream(null);
    } catch (err) {
      setError(err as Error);
      throw err;
    }
  }, [sdk]);

  useEffect(() => {
    if (!sdk) {
      return;
    }

    const handleState = (nextState: RemoteAssistState) => {
      setState(nextState);
    };

    const handleAnswer = (answer: RTCSessionDescriptionInit) => {
      void acceptRemoteAnswer(answer).catch(() => undefined);
    };

    const handleCandidate = (candidate: RTCIceCandidateInit) => {
      void addRemoteIce(candidate).catch(() => undefined);
    };

    const handleTrack = (event: RTCTrackEvent) => {
      const [stream] = event.streams;
      if (stream) {
        setRemoteStream(stream);
      }
    };

    const handleError = (err: Error) => {
      setError(err);
    };

    sdk.on('webrtc:state', handleState);
    sdk.on('webrtc:answer', handleAnswer);
    sdk.on('webrtc:candidate', handleCandidate);
    sdk.on('webrtc:track', handleTrack);
    sdk.on('error', handleError);

    return () => {
      sdk.off('webrtc:state', handleState);
      sdk.off('webrtc:answer', handleAnswer);
      sdk.off('webrtc:candidate', handleCandidate);
      sdk.off('webrtc:track', handleTrack);
      sdk.off('error', handleError);
    };
  }, [sdk, acceptRemoteAnswer, addRemoteIce]);

  return {
    state,
    isActive: state === 'starting' || state === 'offered' || state === 'connecting' || state === 'connected',
    error,
    remoteStream,
    startRemoteAssist,
    acceptRemoteAnswer,
    addRemoteIce,
    endRemoteAssist,
  };
}
