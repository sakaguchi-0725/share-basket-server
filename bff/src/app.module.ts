import { Module } from '@nestjs/common';
import { GraphQLModule } from '@nestjs/graphql';
import { ApolloDriver, ApolloDriverConfig } from '@nestjs/apollo';
import { join } from 'path';
import { PersonalShoppingModule } from './personal-shopping/personal-shopping.module';

@Module({
  imports: [
    GraphQLModule.forRoot<ApolloDriverConfig>({
      driver: ApolloDriver,
      typePaths: ['./**/*.graphql'],
      definitions: {
        path: join(process.cwd(), 'src/graphql/schema.ts'),
        outputAs: 'class',
      },
      playground: true,
    }),
    PersonalShoppingModule,
  ],
})
export class AppModule {}
