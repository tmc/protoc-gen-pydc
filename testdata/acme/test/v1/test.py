# -*- coding: utf-8 -*-
"""
Python Dataclasses for acme.test.v1
"""
from dataclasses import dataclass
from collections import OrderedDict
from enum import Enum
from typing import Dict


class Status(Enum):
    Status_STATUS_UNSPECIFIED = 0 
    Status_STATUS_ACTIVE = 1 
    Status_STATUS_INACTIVE = 2 


@dataclass
class Address:
    street: str 
    city: str 
    state: str 
    country: str 
    address_line: Address_AddressLine 


@dataclass
class Address_AddressLine:
    line: str 


@dataclass
class TestMessage:
    id: str 
    name: str 
    items: str 
    status: Status # Enumeration field
    address: Address # Nested message
    metadata: Dict[str, str] # Map field
    text_payload: str 
    binary_payload: bytes 


@dataclass
class TestMessage_MetadataEntry:
    key: str 
    value: str 


@dataclass
class GetTestMessageRequest:
    id: str 


@dataclass
class GetTestMessageResponse:
    message: TestMessage 


@dataclass
class CreateTestMessageRequest:
    message: TestMessage 


@dataclass
class CreateTestMessageResponse:
    message: TestMessage 
