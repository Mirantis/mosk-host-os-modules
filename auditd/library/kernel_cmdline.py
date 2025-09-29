#!/usr/bin/env python3
# Copyright 2025 Mirantis, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from ansible.module_utils.basic import AnsibleModule

import shlex


ANSIBLE_METADATA = {
    'metadata_version': '0.0.1',
    'status': ['preview'],
    'supported_by': 'dev@mirantis.com'
}

DOCUMENTATION = r'''
---
author:
 - Mirantis (dev@mirantis.com)
module: kernel_cmdline
short_description: Manage kernel cmdline string
version_added: "0.0.1"
description:
  - ''
requirements:
'''


class KernelCmdlineParameters(object):
    def __init__(self, *args):
        self._params = []
        self._changed = False

        for item in args:
            self._params.append(KernelParameter.parse(item))

    @property
    def changed(self):
        return self._changed

    @classmethod
    def parse(cls, string):
        return cls(*shlex.split(string))

    def append(self, name, value=None):
        self._params.append(KernelParameter(name, value))
        self._changed = True

    def set(self, name, value=None, append=False):
        found = False
        for item in self._params:
            if item.name == name:
                found = True
                if item.value != value:
                    item.value = value
                    self._changed = True

        if not found and append:
            self._params.append(KernelParameter(name, value))
            self._changed = True

    def unset(self, name, value=None):
        _params = []
        while self._params:
            item = self._params.pop(0)
            if item.name == name and item.value == value:
                # Skip moving item that has exact name/value
                self._changed = True
            else:
                # Keep item otherwise
                _params.append(item)

        self._params = _params

    def purge(self, name):
        _params = []
        while self._params:
            item = self._params.pop(0)
            if item.name == name:
                # Skip moving item that has exact name
                self._changed = True
            else:
                # Keep item otherwise
                _params.append(item)

        self._params = _params

    def dedup(self, name):
        _params = []
        found = False
        while self._params:
            item = self._params.pop()
            if item.name == name:
                if found:
                    self._changed = True
                else:
                    found = True
                    _params.insert(0, item)
            else:
                _params.insert(0, item)
        self._params = _params

    def __str__(self):
        return ' '.join(map(str, self._params))


class KernelParameter(object):
    def __init__(self, name, value=None):
        self._name = name
        self._value = None

        if value is not None:
            self._value = str(value)

    @property
    def name(self):
        return self._name

    @property
    def value(self):
        return self._value

    @value.setter
    def value(self, value):
        if self._value is None:
            raise Exception("'{}' is a bool argument".format(self._name))
        self._value = str(value)

    @classmethod
    def parse(cls, string):
        if '=' in string:
            name, value = string.split('=')
        else:
            name = string
            value = None
        return cls(name, value)

    def __str__(self):
        if self._value is None:
            return self._name
        elif ' ' in self._value:
            return '{}="{}"'.format(self._name, self._value)
        else:
            return '{}={}'.format(self._name, self._value)


def run_module():
    module = AnsibleModule(
        argument_spec=dict(
            cmdline=dict(type=str, default=""),
            name=dict(type=str, required=True),
            value=dict(type=str, default=None),
            remove_duplicates=(dict(type=bool, default=False)),
            state=dict(type=str, default="present",
                       choices=['append', 'present', 'absent', 'purge']),
        ),
        supports_check_mode=False,
    )

    state = module.params['state']
    name = module.params['name']
    if not name:
        module.fail_json(msg="'name' should be set")

    cmdline = module.params['cmdline']
    value = module.params['value']

    params = KernelCmdlineParameters.parse(cmdline)
    if state == 'append':
        params.append(name, value=value)
    elif state == 'present':
        params.set(name, value=value, append=True)

        if module.params['remove_duplicates']:
            params.dedup(name)
    elif state == 'absent':
        params.unset(name, value=value)
    elif state == 'purge':
        params.purge(name)

    if params.changed:
        module.exit_json(
            changed=True,
            cmdline=str(params),
        )
    else:
        module.exit_json(
            changed=False,
            cmdline=cmdline,
        )


if __name__ == '__main__':
    run_module()
