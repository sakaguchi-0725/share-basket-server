import { Injectable } from '@nestjs/common';
import { LoginInput, SignUpInput } from 'src/graphql/schema';

@Injectable()
export class AuthService {
  login(this: void, input: LoginInput): string {
    console.log(input);
    return 'dummy-access-token';
  }

  signUp(this: void, input: SignUpInput): boolean {
    console.log(input);
    return true;
  }

  verifyToken(this: void, token: string): string {
    console.log(token);
    return 'dummy-user-id';
  }
}
