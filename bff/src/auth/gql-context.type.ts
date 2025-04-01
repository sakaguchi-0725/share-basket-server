import { Request, Response } from 'express';

export type GqlContext = {
  req: Request;
  res: Response;
  userId?: string;
};
