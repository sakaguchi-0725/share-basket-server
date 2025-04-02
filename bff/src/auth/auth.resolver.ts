import { Args, Context, Mutation, Resolver } from '@nestjs/graphql';
import { AuthService } from './auth.service';
import { LoginInput, SignUpInput } from 'src/graphql/schema';
import { Response } from 'express';

@Resolver()
export class AuthResolver {
  constructor(private readonly service: AuthService) {}

  @Mutation(() => Boolean)
  login(
    @Args('input') input: LoginInput,
    @Context('res') res: Response,
  ): boolean {
    const token = this.service.login(input);

    res.cookie('access_token', token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 1000 * 60 * 60 * 24, // 1日
    });

    return true;
  }

  @Mutation(() => Boolean)
  signUp(@Args('input') input: SignUpInput): boolean {
    return this.service.signUp(input);
  }
}
