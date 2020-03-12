===============
Parser Tool Kit
===============

.. image:: https://img.shields.io/github/tag/hydralang/ptk.svg
    :target: https://github.com/hydralang/ptk/tags
.. image:: https://img.shields.io/hexpm/l/plug.svg
    :target: https://github.com/hydralang/ptk/blob/master/LICENSE
.. image:: https://travis-ci.org/hydralang/ptk.svg?branch=master
    :target: https://travis-ci.org/hydralang/ptk
.. image:: https://coveralls.io/repos/github/hydralang/ptk/badge.svg?branch=master
    :target: https://coveralls.io/github/hydralang/ptk?branch=master
.. image:: https://godoc.org/github.com/hydralang/ptk?status.svg
    :target: http://godoc.org/github.com/hydralang/ptk
.. image:: https://img.shields.io/github/issues/hydralang/ptk.svg
    :target: https://github.com/hydralang/ptk/issues
.. image:: https://img.shields.io/github/issues-pr/hydralang/ptk.svg
    :target: https://github.com/hydralang/ptk/pulls
.. image:: https://goreportcard.com/badge/github.com/hydralang/ptk
    :target: https://goreportcard.com/report/github.com/hydralang/ptk

This repository contains a framework for building parsers using the
Pratt technique.  The Pratt technique is a recursive descent parsing
technique specifically geared toward enabling more efficient parsing
of grammars involving operator precedence, such as expression
grammars.  The repository also contains other functionality to assist
in the creation of lexical analyzers and other code useful for
parsers.
