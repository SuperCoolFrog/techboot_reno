:- module(trealla_compat, []).

% 1. Force SWI-Prolog to adhere strictly to the ISO standard
:- set_prolog_flag(iso, true).

% 2. Crash on undefined predicates instead of silently failing or auto-loading
:- set_prolog_flag(unknown, error).

% 3. Mask non-standard SWI built-ins so they throw an error locally
current_predicate(_, _) :-
    throw(error(existence_error(procedure, current_predicate/2), current_predicate/2)).

char_type(_, _) :-
    throw(error(existence_error(procedure, char_type/2), char_type/2)).

% Add any other SWI-specific predicates you want to catch here:
% string/1, split_string/4, etc.

