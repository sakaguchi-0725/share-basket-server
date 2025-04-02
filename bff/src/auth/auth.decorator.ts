import { createParamDecorator, ExecutionContext } from '@nestjs/common';
import { GqlExecutionContext } from '@nestjs/graphql';
import { GqlContext } from './gql-context.type';

export const UserId = createParamDecorator(
  (_, context: ExecutionContext): string | undefined => {
    const ctx = GqlExecutionContext.create(context).getContext<GqlContext>();
    return ctx.userId;
  },
);
