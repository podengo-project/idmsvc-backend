#!/usr/bin/env python3
import sys
import argparse

# Generate an Identity content for the X-Rh-Identity header
# See: https://github.com/RedHatInsights/identity/blob/main/identity.go

class ServiceDetails:
    def __init__(self, is_entitled: bool, is_trial: bool):
        self.is_entitled = is_entitled
        self.is_trial = is_trial

class Entitlements:
    def __init__(self):
        self.data = dict()

    def add(self, name: str, details: ServiceDetails):
        self.data[name] = details
        return self

class UserIdentity:
    def __init__(self):
        super.__init__(self)
        self.username = 'jdoe'
        self.email = 'jdoe@example.com'
        self.first_name = 'John'
        self.last_name = 'Doe'
        self.is_active = True
        self.is_admin = True
        self.is_internal = False
        self.Locale = 'en_EN'
        self.user_id = '12345'

class CertificateIdentity:
    def __init__(self):
        self.subject_dn = '.com.redhat.console'
        self.issuer_dn = '.com.redhat'

class Identity:
    def __init__(self, internal: Internal, user: UserIdentity, associate: AssociateIdentity, certificate: CertificateIdentity, identity_type: str):
        self.identity_type = identity_type
        self.account_number = "12345"
        self.org_id = "12345"
        self.internal = internal
        self.user = user
        self.system = dict()
        self.associate = associate
        self.x509 = certificate

    def build(self):
        if self.type == 'Associate':
            print('%s', self.__build_associate())
        # TODO Check if 'Certificate' is the right string
        elif self.type == 'User':
            print('%s', self.__build_user())
        # TODO Check if 'Certificate' is the right string
        elif self.type == 'Certificate':
            print('%s', self.__build_certificate())
        else:
            sys.exit(1)

    def __build_user(self):
        print('''{
    "identity": %s,
    "entitlements": %s
}''',self.idetity.build(), self.entitlements.build())



class XRHId:
    def __init__(self, identity: Identity, entitlements: Entitlements):
        self.identity = identity
        self.entitlements = entitlements

def parse(args):
    parser = argparse.ArgumentParser(
                prog = 'x-rh-identity.py',
                description = 'Generate a x-rh-identity header to be used with curl',
                epilog = 'Text at the bottom of help')
    parser.add_argument('-t', '--type',
        action='store_const',
        const='user')
    parser.add_argument('-o', '--org-id',
        action='store_const',
        const='12345')
    parser.add_argument('-a', '--account',
        action='store_const',
        const='12345')
    parser.add_argument('--username',
        action='store_const',
        const='jdoe')
    parser.add_argument('--first-name',
        action='store_const',
        const='Jhon')
    parser.add_argument('--last-name',
        action='store_const',
        const='Doe')
    parser.add_argument('--email',
        action='store_const',
        const='jdoe@example.com')


def run(args: list[str]):
    options = parse(args)
    container = XRHId(
        Identity(),
        Entitlements()
    )

if __name__ == "__main__":
    run(sys.argv)
