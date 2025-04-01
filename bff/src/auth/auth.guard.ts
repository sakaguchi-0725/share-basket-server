import {
  CanActivate,
  ExecutionContext,
  Injectable,
  UnauthorizedException,
} from '@nestjs/common';
import { AuthService } from './auth.service';
import { GqlExecutionContext } from '@nestjs/graphql';
import { GqlContext } from './gql-context.type';

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(private readonly service: AuthService) {}

  canActivate(context: ExecutionContext): boolean {
    const ctx = GqlExecutionContext.create(context).getContext<GqlContext>();

    const cookies = ctx.req.cookies as Record<string, string>;
    if (!cookies || !cookies.access_token) {
      throw new UnauthorizedException('アクセストークンが見つかりません');
    }

    const token = cookies.access_token;

    const userId = this.service.verifyToken(token);
    if (!userId) {
      throw new UnauthorizedException('無効なアクセストークンです');
    }

    ctx.userId = userId;
    return true;
  }
}
