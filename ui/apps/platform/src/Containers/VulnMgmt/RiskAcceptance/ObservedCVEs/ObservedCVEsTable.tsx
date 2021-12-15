/* eslint-disable no-nested-ternary */
/* eslint-disable react/no-array-index-key */
import React, { ReactElement, useState } from 'react';
import { TableComposable, Thead, Tbody, Tr, Th } from '@patternfly/react-table';
import {
    Button,
    ButtonVariant,
    Divider,
    DropdownItem,
    InputGroup,
    Pagination,
    TextInput,
    Toolbar,
    ToolbarContent,
    ToolbarItem,
} from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';

import useTableSelection from 'hooks/useTableSelection';

import BulkActionsDropdown from 'Components/PatternFly/BulkActionsDropdown';
import { UsePaginationResult } from 'hooks/patternfly/usePagination';
import DeferralFormModal from './DeferralFormModal';
import FalsePositiveRequestModal from './FalsePositiveFormModal';
import { Vulnerability } from '../imageVulnerabilities.graphql';
import useDeferVulnerability from './useDeferVulnerability';
import useMarkFalsePositive from './useMarkFalsePositive';
import ObservedCVEsTableRow from './ObservedCVEsTableRow';

export type CVEsToBeAssessed = {
    type: 'DEFERRAL' | 'FALSE_POSITIVE';
    ids: string[];
} | null;

export type ObservedCVERow = Vulnerability;

export type ObservedCVEsTableProps = {
    rows: ObservedCVERow[];
    isLoading: boolean;
    registry: string;
    remote: string;
    tag: string;
    itemCount: number;
} & UsePaginationResult;

function ObservedCVEsTable({
    rows,
    registry,
    remote,
    tag,
    itemCount,
    page,
    perPage,
    onSetPage,
    onPerPageSelect,
}: ObservedCVEsTableProps): ReactElement {
    const {
        selected,
        allRowsSelected,
        numSelected,
        onSelect,
        onSelectAll,
        getSelectedIds,
        onClearAll,
    } = useTableSelection<ObservedCVERow>(rows);
    const [cvesToBeAssessed, setCVEsToBeAssessed] = useState<CVEsToBeAssessed>(null);
    const requestDeferral = useDeferVulnerability({
        cveIDs: cvesToBeAssessed?.ids || [],
        registry,
        remote,
        tag,
    });
    const requestFalsePositive = useMarkFalsePositive({
        cveIDs: cvesToBeAssessed?.ids || [],
        registry,
        remote,
        tag,
    });

    function setSelectedCVEsToBeDeferred() {
        const selectedIds = getSelectedIds();
        setCVEsToBeAssessed({ type: 'DEFERRAL', ids: selectedIds });
    }

    function setSelectedCVEsToBeMarkedFalsePositive() {
        const selectedIds = getSelectedIds();
        setCVEsToBeAssessed({ type: 'FALSE_POSITIVE', ids: selectedIds });
    }

    function cancelAssessment() {
        setCVEsToBeAssessed(null);
    }

    function completeAssessment() {
        onClearAll();
        setCVEsToBeAssessed(null);
    }

    return (
        <>
            <Toolbar id="toolbar">
                <ToolbarContent>
                    <ToolbarItem>
                        {/* @TODO: This is just a place holder. Put the correct search filter here */}
                        <InputGroup>
                            <TextInput
                                name="textInput1"
                                id="textInput1"
                                type="search"
                                aria-label="search input example"
                            />
                            <Button
                                variant={ButtonVariant.control}
                                aria-label="search button for search input"
                            >
                                <SearchIcon />
                            </Button>
                        </InputGroup>
                    </ToolbarItem>
                    <ToolbarItem variant="separator" />
                    <ToolbarItem>
                        <BulkActionsDropdown isDisabled={numSelected === 0}>
                            <DropdownItem
                                key="upgrade"
                                component="button"
                                onClick={setSelectedCVEsToBeDeferred}
                            >
                                Defer CVE ({numSelected})
                            </DropdownItem>
                            <DropdownItem
                                key="delete"
                                component="button"
                                onClick={setSelectedCVEsToBeMarkedFalsePositive}
                            >
                                Mark false positive ({numSelected})
                            </DropdownItem>
                        </BulkActionsDropdown>
                    </ToolbarItem>
                    <ToolbarItem variant="pagination" alignment={{ default: 'alignRight' }}>
                        <Pagination
                            itemCount={itemCount}
                            page={page}
                            onSetPage={onSetPage}
                            perPage={perPage}
                            onPerPageSelect={onPerPageSelect}
                        />
                    </ToolbarItem>
                </ToolbarContent>
            </Toolbar>
            <Divider component="div" />
            <TableComposable aria-label="Observed CVEs Table" variant="compact" borders>
                <Thead>
                    <Tr>
                        <Th
                            select={{
                                onSelect: onSelectAll,
                                isSelected: allRowsSelected,
                            }}
                        />
                        <Th>CVE</Th>
                        <Th>Fixable</Th>
                        <Th>Severity</Th>
                        <Th>CVSS score</Th>
                        <Th>Affected components</Th>
                        <Th>Discovered</Th>
                        <Th>Request State</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {rows.map((row, rowIndex) => {
                        const actions = [
                            {
                                title: 'Defer CVE',
                                onClick: (event) => {
                                    event.preventDefault();
                                    setCVEsToBeAssessed({ type: 'DEFERRAL', ids: [row.cve] });
                                },
                            },
                            {
                                title: 'Mark as False Positive',
                                onClick: (event) => {
                                    event.preventDefault();
                                    setCVEsToBeAssessed({ type: 'FALSE_POSITIVE', ids: [row.cve] });
                                },
                            },
                        ];
                        return (
                            <ObservedCVEsTableRow
                                row={row}
                                rowIndex={rowIndex}
                                onSelect={onSelect}
                                selected={selected}
                                actions={actions}
                                page={page}
                                perPage={perPage}
                            />
                        );
                    })}
                </Tbody>
            </TableComposable>
            <DeferralFormModal
                isOpen={cvesToBeAssessed?.type === 'DEFERRAL'}
                numCVEsToBeAssessed={cvesToBeAssessed?.ids.length || 0}
                onSendRequest={requestDeferral}
                onCompleteRequest={completeAssessment}
                onCancelDeferral={cancelAssessment}
            />
            <FalsePositiveRequestModal
                isOpen={cvesToBeAssessed?.type === 'FALSE_POSITIVE'}
                numCVEsToBeAssessed={cvesToBeAssessed?.ids.length || 0}
                onSendRequest={requestFalsePositive}
                onCompleteRequest={completeAssessment}
                onCancelFalsePositive={cancelAssessment}
            />
        </>
    );
}

export default ObservedCVEsTable;